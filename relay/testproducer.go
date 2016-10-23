package main
 
import (
    "flag"
    "net"
    "time"
    "strconv"
    "log"
)
 
func main() {
    var address string
    flag.StringVar(&address, "a", "ec2-54-166-16-218.compute-1.amazonaws.com:8080", "Address of server to connect to")
    flag.Parse()

    conn, err := net.Dial("tcp", address)
    if err  != nil {
        log.Fatal(err)
    } 
    defer conn.Close()

    i := 0
    for {
        msg := strconv.Itoa(i)
        i++
        buf := []byte(msg)
        _,err := conn.Write(buf)
        if err != nil {
            log.Fatal(err)
        }
        if i % 100 == 0 {
            log.Printf("Sent %d bytes", i)
        }
        time.Sleep(time.Nanosecond * 10000000)
    }
}
