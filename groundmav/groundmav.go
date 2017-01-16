package groundmav

import (
  "log"
  "net"
  "fmt"
)

func main() {
  host := "ec2-54-225-61-90.compute-1.amazonaws.com"
  remotePort := 5010
  localPort := 8100

  localConn := connect("", localPort)

  remoteConn := connect(host, remotePort)
  initializeRemoteConnection(remoteConn)

  log.Printf("Bidirectionally routing packets between %s and %s", remoteConn.RemoteAddr(), localConn.RemoteAddr())
  incomingPackets := make(chan []byte, 1000)

  go read(remoteConn, incomingPackets)
  go write(localConn, incomingPackets)

  outgoingPackets := make(chan []byte, 1000)

  go read(localConn, outgoingPackets)
  write(remoteConn, outgoingPackets)
}

func initializeRemoteConnection(c net.Conn) {
  _, err := c.Write([]byte("initiate"))
  handleError(err)
  log.Printf("Sent initializer packet to %s", c.RemoteAddr())
}

func connect(host string, port int) net.Conn {
  destAddr := fmt.Sprintf("%s:%d", host, port)
  destConn, err := net.Dial("udp", destAddr)
  handleError(err)
  log.Printf("Connected to address %s", destAddr)
  return destConn
}

func read(conn net.Conn, packets chan []byte) {
  defer conn.Close()
  for {
    buffer := make([]byte, 150000)
    numRead, err := conn.Read(buffer)
    if err != nil {
      log.Fatal(err)
    }
    packets <- buffer[:numRead]
  }
}

func write(conn net.Conn, packets chan []byte) {
  defer conn.Close()
  for {
    packet := <- packets
    log.Printf("Routed %d bytes", len(packet))
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
