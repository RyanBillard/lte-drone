package main

import (
  "log"
  "net"
)

func main() {
  address := ":8080"
  l, err := net.Listen("tcp", address)
  if err != nil {
    log.Fatal(err)
  }
  defer l.Close()

  connection, err := l.Accept()
  if err != nil {
    log.Fatal(err)
  }
  log.Print("accepted connection from: ", connection.RemoteAddr().String())

  log.Print("attempting to initiate connection to client")

  conn, err := net.Dial("tcp", connection.RemoteAddr().String())
  if err  != nil {
      log.Fatal(err)
  }
  log.Print("successfully connected to client")
  defer conn.Close()
}
