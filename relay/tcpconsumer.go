package main

import (
  "net"
  "fmt"
  "flag"
  "log"
)

func main() {
  var address string
  flag.StringVar(&address, "a", "ec2-54-166-16-218.compute-1.amazonaws.com:5050", "Address of server to connect to")
  flag.Parse()

  conn, err := net.Dial("tcp", address)
  if err != nil {
    log.Fatal(err)
  }
  defer conn.Close()

  for {
    buffer := make([]byte, 1500)
    numRead, err := conn.Read(buffer)
    if err != nil {
      log.Fatal(err)
    }
    fmt.Printf("%s", buffer[:numRead])
  }
}
