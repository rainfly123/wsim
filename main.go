package main

import (
	"log"
	"net/http"
	"encoding/json"
	"runtime"
	"io"
	"./chat"
)

type JsonResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func refreshGrp(w http.ResponseWriter, req *http.Request) {
	groupid := req.FormValue("groupid")

	if len(groupid) < 2 {
		jsonres := JsonResponse{1, "argument error"}
		b, _ := json.Marshal(jsonres)
		io.WriteString(w, string(b))
		return
	}

	chat.Mutex.Lock()
	delete(chat.Groups, groupid)
	chat.Mutex.Unlock()

        jsonres := JsonResponse{1, "OK"}
	b, _ := json.Marshal(jsonres)
	io.WriteString(w, string(b))
	return
}
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
        http.HandleFunc("/refreshgrp", refreshGrp)

	log.Fatal(http.ListenAndServe(":6060", nil))

}
