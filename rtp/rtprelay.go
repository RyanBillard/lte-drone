package main

import (
  "log"
  "net"
)

func main() {
  listenAddr, err := net.ResolveUDPAddr("udp", ":5000")
  clientConn, err := net.ListenUDP("udp", listenAddr)
  handleError(err)
  defer clientConn.Close()

  buffer := make([]byte, 150000)
  numRead, addr, err := clientConn.ReadFromUDP(buffer)
  log.Printf("Received: %v from %v", buffer[:numRead], addr)

  _, err = clientConn.WriteToUDP([]byte("initiate"), addr)
  handleError(err)
  log.Printf("Sent initializer packet to client")
}

func handleError(err error) {
  if err != nil {
    log.Fatal(err)
  }
}
