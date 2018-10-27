// go run ./cmd/kwaygen/main.go -n 10 -m 10 -limit 6

package main

import (
  "flag"
  "fmt"
  "math/rand"
  "os"
  "strconv"
)

var (
  n = flag.Int("n", 5, "number of files")
  m = flag.Int("m", 10, "number of numbers per file")
  limit = flag.Int("limit", 50, "limit of parallel execution")
)

type Result struct {
  Err error
  File string
}

func main() {
  flag.Parse()

  ch := make(chan Result, *n)
  sem := make(chan struct{}, *limit)

  for i := 0; i < *n; i++ {
    sem <- struct{}{}
    go genFile(sem, ch, i, *m)
  }
  for i := 0; i < *n; i++ {
    if r := <- ch; r.Err != nil {
      fmt.Printf("create file error: %v\n", r.Err)
    } else {
      fmt.Printf("created file %q\n", r.File)
    }
  }
}

func genFile(sem chan struct{}, ch chan Result, i, m int) {
  defer func() { <- sem }()

  name := fmt.Sprintf("file.%d", i)
  file, err := os.Create(name)
  if err != nil {
    ch <- Result{
      Err: err,
    }
    return
  }
  defer file.Close()

  for j := 0; j < m; j++ {
    x := rand.Intn(1000)
    s := strconv.Itoa(x)
    _, err := file.WriteString(s + "\n")
    if err != nil {
      ch <- Result{
        Err: err,
      }
      return
    }
  }

  ch <- Result{
    File: name,
  }
}