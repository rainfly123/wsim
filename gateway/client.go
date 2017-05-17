package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if num := len(os.Args); num < 2 {
		fmt.Println("no bindip")
		os.Exit(-1)
	}
	udp_addr, _ := net.ResolveUDPAddr("udp", "192.168.1.240:1111")
	conn, _ := net.ListenUDP("udp", udp_addr)

	dudp_addr, _ := net.ResolveUDPAddr("udp", "192.168.1.240:2222")
	conn.WriteToUDP([]byte("Hello world!"), dudp_addr)

	var msg [20]byte
	conn.Read(msg[0:])

	fmt.Println("msg is", string(msg[0:10]))
}
