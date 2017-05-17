package gateway

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

type Transmit struct {
	RemoteIPAddr string
	Userid       string
	Message      string
}

const LocalIPAddr = "192.168.1.240:1111"

var udp_addr *net.UDPAddr
var conn *net.UDPConn
var ReceiveChan chan Transmit

func SendOut(temp Transmit, remote string) {

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
	udp_addr, _ = net.ResolveUDPAddr("udp", LocalIPAddr)
	conn, _ = net.ListenUDP("udp", udp_addr)
	ReceiveChan = make(chan Transmit, 200)
	fmt.Println("UDPSend inited...")
	go func() {
		for temp := range ReceiveChan {
			SendOut(temp, temp.RemoteIPAddr)
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
