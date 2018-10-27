package main

import (
  "log"
  "net"
  "time"
)

func main() {
  ln, err := net.Listen("tcp", "127.0.0.1:")
  if err != nil {
    panic(err)
  }

  addr := ln.Addr().String()
  log.Printf("listening on %s", addr)

  go func() {
    // Listen here
    conn, err := ln.Accept()
    if err != nil {
      panic(err)
    }
    log.Printf("accepted connection: %s -> %s", conn.LocalAddr(), conn.RemoteAddr())
  }()

  go dial(0, addr)
  go dial(1, addr)

  time.Sleep(time.Hour)
}

func dial(index int, addr string) {
  // Dial here
  log.Printf("#%d dialing %s", index, addr)
  conn, err := net.Dial("tcp", addr)
  if err != nil {
    panic(err)
  }
  log.Printf(
    "#%d successfully dialed: %s -> %s",
    index, conn.LocalAddr(), conn.RemoteAddr(),
  )
}
