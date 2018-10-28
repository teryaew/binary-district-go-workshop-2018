// GOMAXPROCS=1 go run ./cmd/blocksched/main.go

package main

import (
  "fmt"
  "runtime"
  "time"
)

func main() {
  go func() {
    lastTick := time.Now()
    for now := range time.Tick(time.Second) {
      realSleep := now.Sub(lastTick)
      fmt.Printf("hearbeat: %s\n", realSleep)
      lastTick = now

      if diff := realSleep - time.Second; diff > 0 && diff > time.Second {
        panic("WTF")
      }
    }
  }() // Healthchecker goroutine

  go func() {
    for i := 0; ; i++ {
      if i % 500000000000 == 0 {
        runtime.Gosched()
      }
    }
  }()  // Non-cooperative goroutine

  <-time.After(15 * time.Second)
}
