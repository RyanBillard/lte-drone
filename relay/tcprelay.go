package main

import (
  // "fmt"
  "log"
  "net"
)

func main() {

  // Concurrently listen for connections from producer and consumer
  producerChan := make(chan net.Conn)
  go connect(":8080", producerChan)

  consumerChan := make(chan net.Conn)
  go connect(":5050", consumerChan)

  producer := <- producerChan

  consumer := <- consumerChan

  packets := make(chan []byte)
  go write(consumer, packets)
  read(producer, packets)
}

func connect(address string, conn chan net.Conn) {
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
  conn <- connection
}

func read(conn net.Conn, packets chan []byte) {
  defer conn.Close()
  for {
    buffer := make([]byte, 1500)
    _, err := conn.Read(buffer)
    if err != nil {
      log.Fatal(err)
    }
    packets <- buffer
  }
}

func write(conn net.Conn, packets chan []byte) {
  defer conn.Close()
  for {
    packet := <- packets
    _, err := conn.Write(packet)
    if err != nil {
      log.Fatal(err)
    }
  }
}
