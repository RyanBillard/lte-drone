package main

import (
  "log"
  "github.com/RyanBillard/lte-drone/shared"
)

func main() {
  host := "ec2-54-225-61-90.compute-1.amazonaws.com"
  remotePort := 5010
  localPort := 8100

  localConn := shared.Connect("", localPort)

  remoteConn := shared.Connect(host, remotePort)
  shared.InitializeRemoteConnection(remoteConn)

  log.Printf("Bidirectionally routing packets between %s and %s", remoteConn.RemoteAddr(), localConn.RemoteAddr())
  incomingPackets := make(chan []byte, 1000)

  go shared.Read(remoteConn, incomingPackets)
  go shared.Write(localConn, incomingPackets)

  outgoingPackets := make(chan []byte, 1000)

  go shared.Read(localConn, outgoingPackets)
  shared.Write(remoteConn, outgoingPackets)
}
