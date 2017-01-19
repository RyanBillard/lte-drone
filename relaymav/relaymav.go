package main

import (
  "log"
  "net"
  "github.com/RyanBillard/lte-drone/shared"
)

type AddressedPacket struct {
  addr *net.UDPAddr
  data []byte
}

var groundAddr *net.UDPAddr
var droneAddr *net.UDPAddr

func main() {
  //listen for MAVLink
  droneConn := shared.ListenOnPort(8010)

  //wait for ground station to initiate
  clientChan := make(chan shared.Client)
  go shared.DiscoverClient(5010, clientChan)
  client := <- clientChan
  groundAddr = client.Addr

  log.Print("Routing mavlink packets between drone and ground station")

  incomingPackets := make(chan AddressedPacket, 1000)
  go write(client.Conn, incomingPackets)
  go read(droneConn, incomingPackets)

  outgoingPackets := make(chan AddressedPacket, 1000)
  go write(droneConn, outgoingPackets)
  read(client.Conn, outgoingPackets)
}

func write(conn *net.UDPConn, packets chan AddressedPacket) {
  defer conn.Close()
  for {
    packet := <- packets
    log.Printf("Writing %d bytes", len(packet.data))
    _, err := conn.WriteToUDP(packet.data, packet.addr)
    shared.HandleIfError(err)
  }
}

func read(conn *net.UDPConn, packets chan AddressedPacket) {
  defer conn.Close()
  for {
    buffer := make([]byte, 150000)
    numRead, addr, err := conn.ReadFromUDP(buffer)
    shared.HandleIfError(err)
    if droneAddr == nil {
      droneAddr = addr
    }
    log.Printf("Read %d bytes", numRead)
    if addr.String() == droneAddr.String() {
      packets <- AddressedPacket{groundAddr, buffer[:numRead]}
    } else {
      packets <- AddressedPacket{droneAddr, buffer[:numRead]}
    }
  }
}
