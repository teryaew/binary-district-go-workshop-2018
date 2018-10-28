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
  for _, s := range flag.Args() {
    u, err := url.ParseRequestURI(s)
    if err != nil {
      log.Fatalf("invalid url %q: %v", s, err)
    }

    resp, err := http.Get(u.String())
    if err != nil {
      log.Printf("error doing request for %q: %v", u, err)
      continue
    }

    dumpResponse(resp)
  }
}

func dumpResponse(resp *http.Response) {
  dump, err := httputil.DumpResponse(resp, true)
  if err != nil {
    log.Printf("can not dump request from %q", resp.Request.URL)
    return
  }
  log.Printf("got response for %q:\n%s\n\n\n", resp.Request.URL, dump)
}
