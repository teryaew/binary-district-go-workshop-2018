// GOMAXPROCS=1 go run ./cmd/blocksched/main.go

package main

import (
  "fmt"
  "runtime"
  "time"
)

func main() {
  go func() {
    for range time.Tick(time.Second) {
      fmt.Printf("hearbeat\n")
    }
  }() // Healthchecker goroutine

  go func() {
    for i := 0; ; i++ {
      if i % 10000000000 == 0 {
        runtime.Gosched()
      }
    }
  }()  // Non-cooperative goroutine

  <-time.After(10 * time.Second)
}
