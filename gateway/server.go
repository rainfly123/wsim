package gateway

import (
	"../chat"
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

func receiveRouteMessage(server *chat.Server) {
	udp_addr, err := net.ResolveUDPAddr("udp", LocalIPRecvAddr)
	checkError(err)

	conn, err := net.ListenUDP("udp", udp_addr)
	defer conn.Close()
	checkError(err)

	//go recvUDPMsg(conn)
	for {
		var buf [512]byte

		n, raddr, err := conn.ReadFromUDP(buf[0:])
		if err != nil {
			return
		}

		//fmt.Println("msg is ", string(buf[0:n]))
		fmt.Println("addr is ", raddr.String())
		temp := bytes.NewBuffer(buf[:n])
		dec := gob.NewDecoder(&*temp)
		var v Transmit
		err = dec.Decode(&v)
		if err == nil {
			fmt.Println(v)
		}
		client, online := server.GetClient(v.Userid, server)
		if online {
			client.Write([]byte(v.Message))
		}
		//_, err = conn.WriteToUDP([]byte("nice to see u"), raddr)
		//checkError(err)

	}

}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: %s", err.Error())
		os.Exit(1)
	}
}

func InitUDPReceive(server *chat.Server) {
	fmt.Println("UDPReceive inited...")
	go receiveRouteMessage(server)
}
