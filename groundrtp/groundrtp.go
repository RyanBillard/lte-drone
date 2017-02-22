package main

import (
  "log"
  "github.com/stephanernst/lte-drone/shared"
)

func main() {
  rtpHost := "ec2-54-225-61-90.compute-1.amazonaws.com"
  rtpPort := 5000
  rtcpPort := rtpPort + 1

  destRtpPort := 5002
  destRtcpPort := destRtpPort + 1

  destRtpConn := shared.Connect("", destRtpPort)
  destRtcpConn := shared.Connect("", destRtcpPort)

  rtpConn := shared.Connect(rtpHost, rtpPort)
  shared.InitializeRemoteConnection(rtpConn)

  rtcpConn := shared.Connect(rtpHost, rtcpPort)
  shared.InitializeRemoteConnection(rtcpConn)

  log.Printf("Routing rtp stream from server to local ports %d, %d", destRtpPort, destRtcpPort)
  rtpPackets := make(chan []byte, 1000)
  rtcpPackets := make(chan []byte, 1000)

  go shared.Read(rtpConn, rtpPackets)
  go shared.Write(destRtpConn, rtpPackets)

  go shared.Read(rtcpConn, rtcpPackets)
  shared.Write(destRtcpConn, rtcpPackets)
}

