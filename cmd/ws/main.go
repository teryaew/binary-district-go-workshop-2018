package main

import (
  "flag"
  "log"
  "net"
  "sync"

  "github.com/gobwas/ws"
)

var (
  addr = flag.String("addr", "127.0.0.1:8081", "addr to bind to")
)

type Chat struct{
  mu sync.Mutex
  users map[string]*User
}

func NewChat() *Chat {
   return &Chat{
     users: make(map[string]*User),
   }
}

func (c *Chat) Broadcast(msg []byte) {
  c.mu.Lock()
  defer c.mu.Unlock()

  for _, user := range c.users {
    user.Send(msg)
  }
}

func (c *Chat) Register(u *User) {
  c.users[u.Name()] = u
}

func (c *Chat) Remove (u *User) {
  c.mu.Lock()
  defer c.mu.Unlock()
  delete(c.users, u.Name())
}

type User struct{
  conn net.Conn
  sendq chan []byte
}

func NewUser(qsize int, conn net.Conn) *User {
  return &User{
    conn: conn,
    sendq: make(chan []byte, qsize),
  }
}

func (u *User) Name() string {
  return u.conn.RemoteAddr().String()
}

func (u *User) Send(msg []byte) error {
  // FIXME: handle timeouts
  u.sendq <- msg
  return nil
}

func (u *User) drainSendQ() {
  for msg := range u.sendq {
    // FIXME: handle errors
    ws.WriteFrame(u.conn, ws.NewTextFrame(msg))
    u.conn.Write(msg)
  }
}

func (u *User) Recv() ([]byte, error) {
  // Read websocket frame
  f, err := ws.ReadFrame(u.conn)
  if err != nil {
    return nil, err
  }
  return f.Payload, nil
}

func (u *User) Close() {
  u.conn.Close()
}


func main() {
  ln, err := net.Listen("tcp", *addr)
  if err != nil {
    log.Fatal(err)
  }

  c := NewChat()

  for {
    conn, err := ln.Accept()
    if err != nil {
      log.Fatal(err)
    }
    u := NewUser(100, conn)
    c.Register(u)

    go func() {
      msg, err := u.Recv()
      if err != nil {
        c.Remove(u)
        u.Close()
        return
      }

      c.Broadcast(msg)
    }()

    go u.drainSendQ()
  }
}
