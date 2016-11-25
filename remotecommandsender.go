package main

import (
  "log"
  "net"
  "bufio"
  "os"
)

type clientConn struct {
  conn *net.UDPConn
  addr *net.UDPAddr
}

func main() {

  client := waitForClient(":5000")

  packets := make(chan []byte, 1000)
  go write(client.conn, client.addr, packets)
  read(packets)
}

func waitForClient(addr string) clientConn {
  listenAddr, err := net.ResolveUDPAddr("udp", addr)
  handleError(err)
  conn, err := net.ListenUDP("udp", listenAddr)
  handleError(err)

  buffer := make([]byte, 100)
  numRead, clientAddr, err := conn.ReadFromUDP(buffer)
  handleError(err)
  if string(buffer[:numRead]) == "initiate" {
    log.Print("Received initiation message from stream consumer")
  } else {
    log.Fatal("Received unexpected packet from stream consumer")
  }
  return clientConn{conn, clientAddr}
}

func write(conn *net.UDPConn, addr *net.UDPAddr, packets chan []byte) {
  defer conn.Close()
  for {
    packet := <- packets
    log.Print("Writing command")
    _, err := conn.WriteToUDP(packet, addr)
    handleError(err)
  }
}

func read(packets chan []byte) {
  for {
    reader := bufio.NewReader(os.Stdin)
    log.Print("Enter command: ")
    buffer := make([]byte, 1000000)
    numRead, err := reader.Read(buffer)

    handleError(err)
    packets <- buffer[:numRead]
  }
}

func handleError(err error) {
  if err != nil {
    log.Fatal(err)
  }
}
