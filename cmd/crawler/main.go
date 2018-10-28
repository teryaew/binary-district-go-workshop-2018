package main

import (
  "bytes"
  "flag"
  "fmt"
  "github.com/teryaew/bd/pool"
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
    results = make(chan Result)
    done = make(chan bool)
  )

  go func() {
    defer func() { close(done) }()
    for r := range results {
      fmt.Printf("got result %v\n", r)
    }
  }()

  p := pool.NewPool(*parallelism)

  for _, s := range flag.Args() {
    u, err := url.ParseRequestURI(s)
    if err != nil {
      log.Fatalf("invalid url %q: %v", s, err)
    }

    p.Exec(func() {
      r := Result{URL: u}
      r.Count, r.Err = countAt(u)
      results <- r
    })
  }

  p.Close()
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

func dumpResponse(resp *http.Response) {
  dump, err := httputil.DumpResponse(resp, false)
  if err != nil {
    log.Printf("can not dump request from %q", resp.Request.URL)
    return
  }
  log.Printf("got response for %q:\n%s\n\n\n", resp.Request.URL, dump)
}
