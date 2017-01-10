package main

import (
	"fmt"
	"log"

	"./websocket"
)

var origin = "http://localhost/"
var url = "ws://localhost:8080/entry"

func main() {
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	mes := []byte("loginin_234567")
	_, err = ws.Write(mes)
	if err != nil {
		log.Fatal(err)
	}

	message := []byte("emotion_123456_1001_unicast")
	_, err = ws.Write(message)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Send: %s\n", message)

	var msg = make([]byte, 512)
	_, err = ws.Read(msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Receive: %s\n", msg)
	_ = ws.Close()
}
