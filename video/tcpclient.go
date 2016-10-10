package main
 
import (
    "fmt"
    "net"
    "time"
    "strconv"
)
 
func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
    }
}
 
func main() {
    ServerAddr,err := net.ResolveTCPAddr("tcp","ec2-54-166-16-218.compute-1.amazonaws.com:8080")
    CheckError(err)
 
    LocalAddr, err := net.ResolveTCPAddr("tcp", ":0")
    CheckError(err)
 
    Conn, err := net.DialTCP("tcp", LocalAddr, ServerAddr)
    CheckError(err)
 
    defer Conn.Close()
    i := 0
    for {
        msg := strconv.Itoa(i)
        i++
        buf := []byte(msg)
        _,err := Conn.Write(buf)
        if err != nil {
            fmt.Println(msg, err)
        }
        time.Sleep(time.Nanosecond * 10000000)
    }
}
