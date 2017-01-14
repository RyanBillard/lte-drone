package main

import (
  "log"
  "net"
)

type clientConn struct {
  conn *net.UDPConn
  addr *net.UDPAddr
}

func main() {
  //listen for MAVLink
  mavConn := listenOnPort(":8010")

  //wait for client to initiate
  mavChan := make(chan clientConn)
  go waitForClient(":5010", mavChan)
  
  clientMav := <- mavChan

  log.Printf("Routing mavlink packets received at %v to %v", mavConn.LocalAddr(), clientMav.addr.String())
  mavPackets := make(chan []byte, 1000)
  go write(clientMav.conn, clientMav.addr, mavPackets)
  read(mavConn, mavPackets)
}

func waitForClient(addr string, channel chan clientConn) {
  conn := listenOnPort(addr)

  buffer := make([]byte, 100)
  numRead, clientAddr, err := conn.ReadFromUDP(buffer)
  handleError(err)
  if string(buffer[:numRead]) == "initiate" {
    log.Print("Received initiation message from stream consumer")
  } else {
    log.Fatal("Received unexpected packet from stream consumer")
  }
  channel <- clientConn{conn, clientAddr}
}

func listenOnPort(addr string) *net.UDPConn {
  udpAddr, err := net.ResolveUDPAddr("udp", addr)
  handleError(err)
  conn, err := net.ListenUDP("udp", udpAddr)
  handleError(err)
  return conn
}

func write(conn *net.UDPConn, addr *net.UDPAddr, packets chan []byte) {
  defer conn.Close()
  for {
    packet := <- packets
    log.Printf("Writing %d bytes", len(packet))
    _, err := conn.WriteToUDP(packet, addr)
    handleError(err)
  }
}

func read(conn net.Conn, packets chan []byte) {
  defer conn.Close()
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
