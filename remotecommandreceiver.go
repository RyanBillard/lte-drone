package main

import (
  "log"
  "net"
  "fmt"
)

func main() {

  rtpHost := "ec2-54-166-16-218.compute-1.amazonaws.com"
  rtpPort := 5000

  rtpConn := initializeConnection(rtpHost, rtpPort)

  read(rtpConn)
}

func initializeConnection(host string, port int) net.Conn {
  addr := fmt.Sprintf("%s:%d", host, port)
  serverConn, err := net.Dial("udp", addr)
  handleError(err)

  _, err = serverConn.Write([]byte("initiate"))
  handleError(err)
  return serverConn
}

func read(conn net.Conn) {
  defer conn.Close()
  for {
    buffer := make([]byte, 150000)
    numRead, err := conn.Read(buffer)
    if err != nil {
      log.Fatal(err)
    }
    log.Print(buffer[:numRead])
  }
}

func handleError(err error) {
  if err != nil {
    log.Fatal(err)
  }
}
