package shared

import (
  "log"
  "net"
  "fmt"
)

type Client struct {
  Conn *net.UDPConn
  Addr *net.UDPAddr
}

func HandleIfError(err error) {
  if err != nil {
    log.Fatal(err)
  }
}

func Connect(host string, port int) net.Conn {
  destAddr := fmt.Sprintf("%s:%d", host, port)
  destConn, err := net.Dial("udp", destAddr)
  HandleIfError(err)
  log.Printf("Connected to address %s", destAddr)
  return destConn
}

func InitializeRemoteConnection(c net.Conn) {
  _, err := c.Write([]byte("initiate"))
  HandleIfError(err)
  log.Printf("Sent initializer packet to %s", c.RemoteAddr())
}

func Read(conn net.Conn, packets chan []byte) {
  defer conn.Close()
  for {
    buffer := make([]byte, 150000)
    numRead, err := conn.Read(buffer)
    HandleIfError(err)
    packets <- buffer[:numRead]
  }
}

func Write(conn net.Conn, packets chan []byte) {
  defer conn.Close()
  for {
    packet := <- packets
    _, err := conn.Write(packet)
    HandleIfError(err)
    log.Printf("Routed %d bytes", len(packet))
  }
}

func ListenOnPort(port int) *net.UDPConn {
  addr := fmt.Sprintf(":%d", port)
  udpAddr, err := net.ResolveUDPAddr("udp", addr)
  HandleIfError(err)
  conn, err := net.ListenUDP("udp", udpAddr)
  HandleIfError(err)
  return conn
}

func DiscoverClient(port int, channel chan Client) {
  conn := ListenOnPort(port)

  buffer := make([]byte, 100)
  numRead, addr, err := conn.ReadFromUDP(buffer)
  HandleIfError(err)
  if string(buffer[:numRead]) == "initiate" {
    log.Print("Received initiation message from stream consumer")
  } else {
    log.Fatal("Received unexpected packet from stream consumer")
  }
  channel <- Client{conn, addr}
}
