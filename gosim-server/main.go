package main

import (
  "log"
)

const SERVER_IP = ":8080"

func main() {
  s := server{
    ip:       SERVER_IP,
    sessions: []*session{},
  }
  err := s.startServer()
  if err != nil {
    log.Fatal(err)
  }
}
