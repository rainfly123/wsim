package main

import (
	"log"
	"net/http"

	"./chat"
)

func main() {
	log.SetFlags(log.Lshortfile)
	chat.InitGroup()

	// websocket server
	server := chat.NewServer("/entry")
	go server.Listen()
	go chat.RecGrpMsgTrd(server)

	// static files
	http.Handle("/", http.FileServer(http.Dir("webroot")))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
