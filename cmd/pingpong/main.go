package main

import (
  "fmt"
  "runtime"
  "time"
)

type Ball struct{}

func main() {
  runtime.GOMAXPROCS(1)
  table := make(chan Ball)
  go player(table, "Ivan")
  go player(table, "Petr")
  table <- Ball{}
  time.Sleep(1 * time.Second)
  <- table // Grab the ball
  close(table)
  fmt.Println("Game over")
  time.Sleep(time.Second)
}

func player(table chan Ball, name string) {
  defer fmt.Printf("%s gone\n", name)
  for ball := range table {
    fmt.Printf("Yay, %s got the ball!\n", name)
    time.Sleep(time.Second)
    table <- ball
  }
}
