package main

import (
  "log"
  "net"
)

func main() {
  rtpAddr, err := net.ResolveUDPAddr("udp", ":8000")
  handleError(err)
  rtpConn, err := net.ListenUDP("udp", rtpAddr)
  handleError(err)
  defer rtpConn.Close()

  listenAddr, err := net.ResolveUDPAddr("udp", ":5000")
  handleError(err)
  clientConn, err := net.ListenUDP("udp", listenAddr)
  handleError(err)
  defer clientConn.Close()

  buffer := make([]byte, 150000)
  numRead, clientAddr, err := clientConn.ReadFromUDP(buffer)
  handleError(err)
  if string(buffer[:numRead]) == "initiate" {
    log.Print("Received initiation message from stream consumer")
  } else {
    log.Fatal("Received unexpected packet from strem consumer")
  }

  log.Printf("Routing rtp messages received on port %v to %v", rtpAddr.String(), clientAddr.String())
  packets := make(chan []byte, 1000)
  go write(clientConn, clientAddr, packets)
  read(rtpConn, packets)
}

func write(conn *net.UDPConn, addr *net.UDPAddr, packets chan []byte) {
  for {
    packet := <- packets
    log.Printf("Writing %d bytes", len(packet))
    _, err := conn.WriteToUDP(packet, addr)
    handleError(err)
  }
}

func read(conn net.Conn, packets chan []byte) {
  for {
    buffer := make([]byte, 150000)
    numRead, err := conn.Read(buffer)
    handleError(err)
    log.Printf("Read %d bytes", numRead)
    packets <- buffer[:numRead]
  }
}

func handleError(err error) {
  if err != nil {
    log.Fatal(err)
  }
}
