// GOMAXPROCS=1 go run ./cmd/blocksched/main.go

package main

import (
  "fmt"
  "time"
)

func main() {
  go func() {
    for range time.Tick(time.Second) {
      fmt.Printf("hearbeat\n")
    }
  }() // Healthchecker goroutine

  go func() {
    for {}
  }()  // Non-cooperative goroutine

  <-time.After(10 * time.Second)
}
