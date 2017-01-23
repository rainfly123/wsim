package main

import (
	"fmt"
	"log"

	"./websocket"
)

var origin = "http://localhost/"
var url = "ws://live.66boss.com:6060/entry"

func main() {
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	mes := []byte("login_1000005730")
	_, err = ws.Write(mes)
	if err != nil {
		log.Fatal(err)
	}
	/*
			var msg = make([]byte, 512)
			_, err = ws.Read(msg)
			if err != nil {
				log.Fatal(err)
			}
		fmt.Printf("Receive: %s\n", msg)
	*/
	var message []byte
	message = make([]byte, 1024)
	for i := 0; i < 10; i++ {
		n, err := ws.Read(message)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(message[:n]))
	}

	//	var msg = make([]byte, 512)
	//	_, err = ws.Read(msg)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Printf("Receive: %s\n", msg)
	_ = ws.Close()
}
