package duckTyping

import (
  "fmt"
  "os"
  "strconv"
)

type Greeter interface {
  Greet() string
}

type A string

func (a A) Greet() string {
  return string(a)
}

type B struct {
  Greeting string
  Count int
}

func (b B) Greet() string {
  var ret string
  for i := 0; i < b.Count; i++ {
    ret += b.Greeting
  }
  //return strings.Repeat(b.Greeting, b.Count)
  return ret
}

func Hello(g Greeter) {
  fmt.Println(g.Greet())
}

func main() {
  if len(os.Args) < 2 {
    fmt.Println("no greeter given")
    os.Exit(1)
  }

  var g Greeter
  switch os.Args[1] {
  case "a":
    g = A("hello there1")
  case "b":
    var count int
    if len(os.Args) >= 3 {
      var err error
      count, err = strconv.Atoi(os.Args[2])
      if err != nil {
        fmt.Printf("error reading count: %v", err)
        os.Exit(1)
      }
    } else {
      count = 8
    }

    g = B{
      Greeting: "hi!",
      Count: count,
    }
  default:
    fmt.Println("no such greeter")
    os.Exit(1)
  }
  Hello(g)
}
