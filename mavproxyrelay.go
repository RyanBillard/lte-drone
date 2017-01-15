package main

import (
  "log"
  "net"
)

type AddressedPacket struct {
  addr *net.UDPAddr
  data []byte
}

func main() {
  //listen for MAVLink
  droneConn := listenOnPort(":8010")

  groundConn := listenOnPort(":5010")
  //wait for ground station to initiate
  waitForClient(groundConn)

  log.Print("Routing mavlink packets between drone and ground station")

  incomingPackets := make(chan AddressedPacket, 1000)
  go write(groundConn, incomingPackets)
  read(droneConn, incomingPackets)
}

func waitForClient(conn *net.UDPConn) {
  buffer := make([]byte, 100)
  numRead, _, err := conn.ReadFrom(buffer)
  handleError(err)
  if string(buffer[:numRead]) == "initiate" {
    log.Print("Received initiation message from stream consumer")
  } else {
    log.Fatal("Received unexpected packet from stream consumer")
  }
}

func listenOnPort(addr string) *net.UDPConn {
  udpAddr, err := net.ResolveUDPAddr("udp", addr)
  handleError(err)
  conn, err := net.ListenUDP("udp", udpAddr)
  handleError(err)
  return conn
}

func write(conn *net.UDPConn, packets chan AddressedPacket) {
  defer conn.Close()
  for {
    packet := <- packets
    log.Printf("Writing %d bytes", len(packet.data))
    _, err := conn.WriteToUDP(packet.data, packet.addr)
    handleError(err)
  }
}

func read(conn *net.UDPConn, packets chan AddressedPacket) {
  defer conn.Close()
  for {
    buffer := make([]byte, 150000)
    numRead, addr, err := conn.ReadFromUDP(buffer)
    handleError(err)
    log.Printf("Read %d bytes", numRead)
    packets <- AddressedPacket{addr, buffer[:numRead]}
  }
}

func handleError(err error) {
  if err != nil {
    log.Fatal(err)
  }
}
