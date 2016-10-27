package main

import (
  "log"
  "net"
  "flag"
)

func main() {
  // ec2-54-166-16-218.compute-1.amazonaws.com:5000
  var listenAddr string
  flag.StringVar(&listenAddr, "listen", ":5000", "Address of server from which rtp stream is requested")

  var destAddr string
  flag.StringVar(&destAddr, "destination", ":5002", "Address to which rtp stream data will be routed")

  flag.Parse()

  serverConn, err := net.Dial("udp", listenAddr)
  handleError(err)
  defer serverConn.Close()
  log.Printf("Connected to server at %v", serverConn.RemoteAddr())

  destUDPAddr, err := net.ResolveUDPAddr("udp", destAddr)
  handleError(err)
  destConn, err := net.DialUDP("udp", nil, destUDPAddr)
  handleError(err)
  defer destConn.Close()
  log.Printf("Connected to dest port at %v", destConn.RemoteAddr())

  _, err = serverConn.Write([]byte("initiate"))
  handleError(err)
  log.Printf("Sent initializer packet to server")

  log.Printf("Routing messages received on port %v to local port %v", listenAddr, destAddr)
  packets := make(chan []byte, 1000)
  go read(serverConn, packets)
  write(destConn, destUDPAddr, packets)
}

func read(conn net.Conn, packets chan []byte) {
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
