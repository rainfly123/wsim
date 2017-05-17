package main

import (
	"fmt"
	"net"
	"os"
)

var bindip string

func main() {
	if num := len(os.Args); num < 2 {
		fmt.Println("no bindip")
		os.Exit(-1)
	}
	fmt.Println("Bind IP Addr:", os.Args[1])

	udp_addr, err := net.ResolveUDPAddr("udp", os.Args[1]+":2222")
	checkError(err)

	conn, err := net.ListenUDP("udp", udp_addr)
	defer conn.Close()
	checkError(err)

	//go recvUDPMsg(conn)
	for {
		recvUDPMsg(conn)
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: %s", err.Error())
		os.Exit(1)
	}
}

func recvUDPMsg(conn *net.UDPConn) {
	var buf [20]byte

	n, raddr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		return
	}

	fmt.Println("msg is ", string(buf[0:n]))
	fmt.Println("addr is ", raddr.String())

	_, err = conn.WriteToUDP([]byte("nice to see u"), raddr)
	checkError(err)
}
