package main

import (
  "log"
  "net"
)

type server struct {
  ip       string
  sessions []*session
}

func (serv *server) startServer() error {
  log.Println("Starting server on port 8080")
  l, err := net.Listen("tcp", serv.ip)
  if err != nil {
    log.Println("Error starting server to listen on port 8080", err)
    return err
  }
  defer l.Close()

  for {
    conn, err := l.Accept()
    if err != nil {
      log.Println("Error accepting connection", err)
      continue
    }
    go func() {
      err := serv.handleConnection(conn)
      if err != nil {
        log.Println("Error handling connection", err)
        return
      }
    }()
  }
}

func (serv *server) handleConnection(conn net.Conn) error {
  defer conn.Close()
  log.Println("Client connected:", conn.RemoteAddr().String())

  // create a new session for the client
  s := &session{
    conn:           conn,
    user:           user{},
    incomingEvents: make(chan event),
    outgoingEvents: make(chan event),
    closed:         false,
  }

  serv.sessions = append(serv.sessions, s)

  // start the incoming event loop
  go func() {
    for s.closed != true {
      err := s.recvIncoming()
      if err != nil {
        log.Fatal(err)
      }
    }
  }()

  // start the event handling loop
  go func() {
    for s.closed != true {
      err := s.handleEvent()
      if err != nil {
        log.Fatal(err)
      }
    }
  }()

  // start the outgoing event loop
  go func() {
    for s.closed != true {
      err := s.sendOutgoing()
      if err != nil {
        log.Fatal(err)
      }
    }
  }()

  // send a welcome message to the client
  s.sendChatMessage("Server", "Welcome to the gosim server!")

  // loop to keep the connection alive until the client disconnects
  for {
    // if err := conn.SetReadDeadline(time.Now().Add(readTimeout)); err != nil {
    //  return err
    // }

    if s.closed == true {
      s.conn.Close()
    }
  }
}
