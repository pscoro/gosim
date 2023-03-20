package main

import (
  "fmt"
  "log"
  "net"
)

type server struct {
  ip      string
  session session
}

func connectToServer(ip string) (*server, error) {
  c, err := net.Dial("tcp", ip)
  if err != nil {
    fmt.Println(err)
    return nil, err
  }

  // create a new session for the client
  session := &session{
    conn:           c,
    user:           user{},
    incomingEvents: make(chan event),
    outgoingEvents: make(chan event),
  }
  server := &server{
    ip:      ip,
    session: *session,
  }
  go func() {
    for {
      err := server.session.recvIncoming()
      if err != nil {
        log.Fatal(err)
      }
    }
  }()
  go func() {
    for {
      err := server.session.handleEvent()
      if err != nil {
        log.Fatal(err)
      }
    }
  }()
  go func() {
    for {
      err := server.session.sendOutgoing()
      if err != nil {
        log.Fatal(err)
      }
    }
  }()
  return server, nil
}
