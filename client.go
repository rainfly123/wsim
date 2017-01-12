package main

import (
	"fmt"
	"log"
	"time"

	"./websocket"
)

var origin = "http://localhost/"
var url = "ws://localhost:8080/entry"

func main() {
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	mes := []byte("login_234567")
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

	for i := 0; i < 1; i++ {
		message := []byte("emotion_5487_1001_group_extension")
		//	message := []byte("emotion_1000005730_1001_unicast_extension")
		_, err = ws.Write(message)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Send: %s\n", message)
		time.Sleep(1 * time.Second)
	}

	//	var msg = make([]byte, 512)
	//	_, err = ws.Read(msg)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Printf("Receive: %s\n", msg)
	_ = ws.Close()
}
