package main

import (
  "log"
  "net"
  "flag"
)

func main() {
  // ec2-54-166-16-218.compute-1.amazonaws.com:5050
  var listenAddr string
  flag.StringVar(&listenAddr, "listen", ":5000", "Address of server from which rtp stream is requested")

  // var destAddr string
  // flag.StringVar(&destAddr, "destination", ":8000", "Address to which rtp stream data will be routed")

  flag.Parse()

  serverConn, err := net.Dial("udp", listenAddr)
  handleError(err)
  defer serverConn.Close()
  log.Printf("Connected to server at %v", serverConn.RemoteAddr())

  _, err = serverConn.Write([]byte("initiate"))
  handleError(err)
  log.Printf("Sent initializer packet")

  buffer := make([]byte, 150000)
  numRead, err := serverConn.Read(buffer)
  log.Printf("Received: %v from %v", buffer[:numRead], serverConn.RemoteAddr())
}

func handleError(err error) {
  if err != nil {
    log.Fatal(err)
  }
}
