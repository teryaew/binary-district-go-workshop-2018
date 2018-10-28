package main

import (
  "bytes"
  "flag"
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "net/http/httputil"
  "net/url"
)

var (
  parallelism = flag.Int("parallelism", 5, "number of parallel requests")
)

type Result struct{
  Err error
  URL *url.URL
  Count int
}

func main() {
  flag.Parse()

  var (
    sem = make(chan struct{}, *parallelism)
    work = make(chan *url.URL)
    results = make(chan Result)
    done = make(chan bool)
  )

  go func() {
    defer func() { close(done) }()
    for r := range results {
      fmt.Printf("got result %v\n", r)
    }
  }()

  for _, s := range flag.Args() {
    u, err := url.ParseRequestURI(s)
    if err != nil {
      log.Fatalf("invalid url %q: %v", s, err)
    }

    select {
    case work <- u:
      default:
        select {
        case sem <- struct{}{}:
          go process(sem, work, results, u)
        case work <- u:
        }
    }
  }

  close(work)

  // Wait for all workers are done
  for i := 0; i < *parallelism; i++ {
    sem <- struct{}{}
  }

  close(results)
  <-done
}

func countAt(u *url.URL) (int, error) {
  resp, err := http.Get(u.String())
  if err != nil {
    return 0, err
  }

  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return 0, err
  }
  return bytes.Count(body, []byte("go")), nil
}

func process(sem <-chan struct{}, work <-chan *url.URL, results chan<- Result, u *url.URL) {
  defer func() { <-sem }()
  var ok bool
  for {
    r := Result{URL: u}
    r.Count, r.Err = countAt(u)
    results <- r

    if u, ok = <-work; !ok {
      return
    }
  }
}

func dumpResponse(resp *http.Response) {
  dump, err := httputil.DumpResponse(resp, false)
  if err != nil {
    log.Printf("can not dump request from %q", resp.Request.URL)
    return
  }
  log.Printf("got response for %q:\n%s\n\n\n", resp.Request.URL, dump)
}
