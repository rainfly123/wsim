package main

import (
	"./chat"
	"encoding/json"
	"io"
	"log"
	"menteslibres.net/gosexy/redis"
	"net/http"
	"runtime"
	"strings"
)

type JsonResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
type JsonResponseData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
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

func createHandle(w http.ResponseWriter, req *http.Request) {
	creator := req.FormValue("creator")
	members := req.FormValue("members")

	if len(creator) < 2 || len(members) < 2 {
		jsonres := JsonResponse{1, "argument error"}
		b, _ := json.Marshal(jsonres)
		io.WriteString(w, string(b))
		return
	}
	var client *redis.Client
	var ok bool

	amembers := strings.TrimSpace(members)
	bmembers := strings.TrimSuffix(amembers, ",")
	users := strings.Split(bmembers, ",")

	client, ok = chat.Clients.Get()
	if ok != true {
		log.Panic("redis error")
		return
	}

	strID, _ := client.Get("groupID")
	key := "group_" + strID
	mkey := "groupmembers_" + strID

	client.HMSet(key, "creator", creator, "name", "", "notice", "", "qrcode", "")
	for _, v := range users {
		client.SAdd(mkey, v)
		skey := "mygroups_" + v
		client.SAdd(skey, strID)
	}

	client.Incr("groupID")
	client.Close()

	jsonres := JsonResponseData{1, "OK", strID}
	b, _ := json.Marshal(jsonres)
	io.WriteString(w, string(b))
	return
}

func addHandle(w http.ResponseWriter, req *http.Request) {
	groupid := req.FormValue("groupid")
	members := req.FormValue("members")

	if len(members) < 2 || len(groupid) < 2 {
		jsonres := JsonResponse{1, "argument error"}
		b, _ := json.Marshal(jsonres)
		io.WriteString(w, string(b))
		return
	}
	var client *redis.Client
	var ok bool

	amembers := strings.TrimSpace(members)
	bmembers := strings.TrimSuffix(amembers, ",")
	users := strings.Split(bmembers, ",")

	client, ok = chat.Clients.Get()
	if ok != true {
		log.Panic("redis error")
		return
	}

	mkey := "groupmembers_" + groupid

	for _, v := range users {
		client.SAdd(mkey, v)
		skey := "mygroups_" + v
		client.SAdd(skey, groupid)
	}

	client.Close()

	jsonres := JsonResponse{1, "OK"}
	b, _ := json.Marshal(jsonres)
	io.WriteString(w, string(b))
	return
}

func delHandle(w http.ResponseWriter, req *http.Request) {
	groupid := req.FormValue("groupid")
	members := req.FormValue("members")

	if len(members) < 2 || len(groupid) < 2 {
		jsonres := JsonResponse{1, "argument error"}
		b, _ := json.Marshal(jsonres)
		io.WriteString(w, string(b))
		return
	}
	var client *redis.Client
	var ok bool

	amembers := strings.TrimSpace(members)
	bmembers := strings.TrimSuffix(amembers, ",")
	users := strings.Split(bmembers, ",")

	client, ok = chat.Clients.Get()
	if ok != true {
		log.Panic("redis error")
		return
	}

	mkey := "groupmembers_" + groupid

	for _, v := range users {
		client.SRem(mkey, v)
		skey := "mygroups_" + v
		client.SRem(skey, groupid)
	}

	client.Close()

	jsonres := JsonResponse{1, "OK"}
	b, _ := json.Marshal(jsonres)
	io.WriteString(w, string(b))
	return
}
func editHandle(w http.ResponseWriter, req *http.Request) {
	groupid := req.FormValue("groupid")
	notice := req.FormValue("notice")
	name := req.FormValue("name")

	if len(groupid) < 2 {
		jsonres := JsonResponse{1, "argument error"}
		b, _ := json.Marshal(jsonres)
		io.WriteString(w, string(b))
		return
	}
	var client *redis.Client
	var ok bool

	client, ok = chat.Clients.Get()
	if ok != true {
		log.Panic("redis error")
		return
	}
	key := "group_" + groupid
	if len(notice) > 2 {
		client.HSet(key, "notice", notice)
	}
	if len(name) > 2 {
		client.HSet(key, "name", name)
	}
	client.Close()

	jsonres := JsonResponse{1, "OK"}
	b, _ := json.Marshal(jsonres)
	io.WriteString(w, string(b))
	return
}
func querymygrpHandle(w http.ResponseWriter, req *http.Request) {
	userid := req.FormValue("userid")

	if len(userid) < 2 {
		jsonres := JsonResponse{1, "argument error"}
		b, _ := json.Marshal(jsonres)
		io.WriteString(w, string(b))
		return
	}
	var client *redis.Client
	var ok bool

	client, ok = chat.Clients.Get()
	if ok != true {
		log.Panic("redis error")
		return
	}
	key := "mygroups_" + userid
	groups, _ := client.SMembers(key)
	client.Close()
	type MyResopnse struct {
		JsonResponse
		Groups []string `json:"data"`
	}
	jsonres := MyResopnse{JsonResponse{1, "OK"}, groups}
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
	http.HandleFunc("/creategrp", createHandle)
	http.HandleFunc("/addmembers", addHandle)
	http.HandleFunc("/delmembers", delHandle)
	http.HandleFunc("/editgrp", editHandle)
	http.HandleFunc("/querymygrps", querymygrpHandle)

	log.Fatal(http.ListenAndServe(":6060", nil))

}
