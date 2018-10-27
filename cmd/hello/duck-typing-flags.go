package duckTyping

import (
  "flag"
  "fmt"
  "os"
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
  var (
    t = flag.String("type", "a", "type of greeting")
    n = flag.Int("repeat", 8, "number of greetings for b")
  )
  flag.Parse()

  var g Greeter
  switch *t {
  case "a":
    g = A("hello there1")
  case "b":
    g = B{
      Greeting: "hi!",
      Count: *n,
    }
  default:
    fmt.Println("no such greeter")
    os.Exit(1)
  }
  Hello(g)
}
