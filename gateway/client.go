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

	conn, err := net.Dial("udp", os.Args[1]+":2222")
	defer conn.Close()
	if err != nil {
		return
	}
	conn.Write([]byte("Hello world!"))

	var msg [20]byte
	conn.Read(msg[0:])

	fmt.Println("msg is", string(msg[0:10]))
}
