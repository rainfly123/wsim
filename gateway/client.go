package gateway

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	//	"time"
)

var udp_addr *net.UDPAddr
var conn *net.UDPConn
var receiveChan chan Transmit

func sendOut(temp Transmit, remote string) {

	dudp_addr, _ := net.ResolveUDPAddr("udp4", remote)

	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err2 := enc.Encode(temp)
	if err2 == nil {
		conn.WriteToUDP(network.Bytes(), dudp_addr)
	}

	//var msg [20]byte
	//conn.Read(msg[0:])
}

func InitUDPSend() {
	udp_addr, _ = net.ResolveUDPAddr("udp", LocalIPSendAddr)
	conn, _ = net.ListenUDP("udp", udp_addr)
	receiveChan = make(chan Transmit, 200)
	fmt.Println("UDPSend inited...")
	go func() {
		for temp := range receiveChan {
			sendOut(temp, temp.RemoteIPAddr)
			fmt.Println(temp)
		}
	}()
}

/*
func main() {
	var temp = Transmit{"192.168.1.240:2222", "100000320", "video_touser_from啊地方"}
	InitUDPSend()
	for {
		ReceiveChan <- temp
		time.Sleep(1 * time.Second)
	}
}
*/
