package main

import (
  "flag"
  "fmt"
  "log"
  "net/http"
  "net/http/httputil"
  "net/url"
)

var (
  parallelism = flag.Int("parallelism", 5, "number of parallel requests")
)

type Result struct{}

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

func process(sem <-chan struct{}, work <-chan *url.URL, results chan<- Result, u *url.URL) {
  defer func() { <-sem }()
  for {
    resp, err := http.Get(u.String())
    if err != nil {
      log.Printf("error doing request for %q: %v", u, err)
      continue
    }
    dumpResponse(resp)

    results <- Result{}

    var ok bool
    u, ok = <-work
    if !ok {
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
