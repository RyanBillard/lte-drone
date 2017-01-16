package main

import (
  "log"
  "net"
  "fmt"
)

func main() {
  rtpHost := "ec2-54-225-61-90.compute-1.amazonaws.com"
  rtpPort := 5000
  rtcpPort := rtpPort + 1

  destRtpPort := 5002
  destRtcpPort := destRtpPort + 1

  destRtpConn, destRtpAddr := initialiazeLocalConnection(destRtpPort)
  destRtcpConn, destRtcpAddr := initialiazeLocalConnection(destRtcpPort)

  rtpConn := initializeConnection(rtpHost, rtpPort)
  rtcpConn := initializeConnection(rtpHost, rtcpPort)

  log.Printf("Routing rtp stream from server to local ports %d, %d", destRtpPort, destRtcpPort)
  rtpPackets := make(chan []byte, 1000)
  rtcpPackets := make(chan []byte, 1000)

  go read(rtpConn, rtpPackets)
  go write(destRtpConn, destRtpAddr, rtpPackets)

  go read(rtcpConn, rtcpPackets)
  write(destRtcpConn, destRtcpAddr, rtcpPackets)
}

func initializeConnection(host string, port int) net.Conn {
  addr := fmt.Sprintf("%s:%d", host, port)
  serverConn, err := net.Dial("udp", addr)
  handleError(err)
  // defer serverConn.Close()
  log.Printf("Connected to server at %v", serverConn.RemoteAddr())

  _, err = serverConn.Write([]byte("initiate"))
  handleError(err)
  log.Printf("Sent initializer packet to server")
  return serverConn
}

func initialiazeLocalConnection(port int) (*net.UDPConn, *net.UDPAddr) {
  destAddr := fmt.Sprintf(":%d", port)
  destUDPAddr, err := net.ResolveUDPAddr("udp", destAddr)
  handleError(err)
  destConn, err := net.DialUDP("udp", nil, destUDPAddr)
  handleError(err)
  // defer destConn.Close()
  log.Printf("Connected to dest port at %v", destConn.RemoteAddr())
  return destConn, destUDPAddr
}

func read(conn net.Conn, packets chan []byte) {
  defer conn.Close()
  for {
    buffer := make([]byte, 150000)
    numRead, err := conn.Read(buffer)
    if err != nil {
      log.Fatal(err)
    }
    log.Printf("Read %d bytes", numRead)
    packets <- buffer[:numRead]
  }
}

func write(conn *net.UDPConn, addr *net.UDPAddr, packets chan []byte) {
  defer conn.Close()
  for {
    packet := <- packets
    log.Printf("Writing %d bytes", len(packet))
    _, err := conn.Write(packet)
    if err != nil {
      log.Fatal(err)
    }
  }
}

func handleError(err error) {
  if err != nil {
    log.Fatal(err)
  }
}
