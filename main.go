package main

import (
	"log"
	"net/http"
	"runtime"

	"./chat"
)

func main() {
	runtime.GOMAXPROCS(4)
	log.SetFlags(log.Lshortfile)
	chat.InitGroup()
	chat.InitRedis()

	// websocket server
	server := chat.NewServer("/entry")
	go server.Listen()
	go chat.RecGrpMsgTrd(server)
	//	go chat.HeaartbeatTrd(server)

	// static files
	http.Handle("/", http.FileServer(http.Dir("webroot")))

	log.Fatal(http.ListenAndServe(":6060", nil))
}
