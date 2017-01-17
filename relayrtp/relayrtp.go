package relayrtp

import (
  "log"
  "net"
  "github.com/RyanBillard/lte-drone/shared"
)

func main() {
  //listen for RTP stream
  rtpConn := shared.ListenOnPort(8000)
  rtcpConn := shared.ListenOnPort(8001)

  //wait for client to initiate
  rtpChan := make(chan shared.Client)
  go shared.DiscoverClient(5000, rtpChan)
  rtcpChan := make(chan shared.Client)
  go shared.DiscoverClient(5001, rtcpChan)
  
  clientRtp := <- rtpChan 
  clientRtcp := <- rtcpChan
  
  log.Printf("Routing rtp messages received on ports %v, %v to %v, %v", rtpConn.LocalAddr(), rtcpConn.LocalAddr(), clientRtp.Addr.String(), clientRtcp.Addr.String())
  rtpPackets := make(chan []byte, 1000)
  go write(clientRtp.Conn, clientRtp.Addr, rtpPackets)
  go shared.Read(rtpConn, rtpPackets)

  rtcpPackets := make(chan []byte, 1000)
  go write(clientRtcp.Conn, clientRtcp.Addr, rtcpPackets)
  shared.Read(rtcpConn, rtcpPackets)
}

func write(conn *net.UDPConn, addr *net.UDPAddr, packets chan []byte) {
  defer conn.Close()
  for {
    packet := <- packets
    log.Printf("Writing %d bytes", len(packet))
    _, err := conn.WriteToUDP(packet, addr)
    shared.HandleIfError(err)
  }
}
