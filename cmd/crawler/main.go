package main

import (
  "flag"
  "log"
  "net/http"
  "net/http/httputil"
  "net/url"
)

var (
  parallelism = flag.Int("parallelism", 5, "number of parallel requests")
)

func main() {
  flag.Parse()

  var (
    sem = make(chan struct{}, *parallelism)
    work = make(chan *url.URL)
  )

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
          go process(work, u)
        case work <- u:
        }
    }
  }
}

func process(work <-chan *url.URL, u *url.URL) {
  for {
    resp, err := http.Get(u.String())
    if err != nil {
      log.Printf("error doing request for %q: %v", u, err)
      continue
    }
    dumpResponse(resp)
    u = <-work
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
