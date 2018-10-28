package main

import "fmt"

func main() {
  x := make([]int, 5, 5)
  y := x[1:3] // x[1:3:3] to increase capacity -> z will be a new array
  z := append(y, 42)

  fmt.Println(x)
  fmt.Println(y)
  fmt.Println(z)
}
