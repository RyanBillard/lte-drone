package main
import (
  "syscall"
  // "os"
  "log"
)

func main() {
  connectSock := createSocket()
  acceptSock := createSocket()

  // remoteAddr := syscall.SockaddrInet4{Port: 8080, Addr: [4]byte{0,0,0,0}}
  remoteAddr := syscall.SockaddrInet4{Port: 8080, Addr: [4]byte{54,166,16,218}}
  err := syscall.Connect(connectSock, &remoteAddr)
  handleError(err, connectSock)
  log.Print("Connected to: ", remoteAddr)

  var n int = 1
  err = syscall.Listen(acceptSock, n)
  handleError(err, acceptSock)

  _, connectedAddr, err := syscall.Accept(acceptSock)
  handleError(err, acceptSock)

  log.Print("Accepted connection from: ", connectedAddr)

  // file := os.NewFile(uintptr(fd), "")
}

func createSocket() int {
    fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
    handleError(err, fd)

    // set REUSEPORT
    err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEPORT, 1)
    handleError(err, fd)

    // set REUSEADDR
    err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
    handleError(err, fd)

    // bind to local port
    localAddr := syscall.SockaddrInet4{Port: 5000, Addr: [4]byte{0,0,0,0}}
    err = syscall.Bind(fd, &localAddr)
    handleError(err, fd)

    return fd
}

func handleError(err error, fd int) {
  if err != nil {
    syscall.Close(fd)
    log.Fatal(err)
  }
}
