package main

import (
	"./chat"
	"./snap"
	"./userinfo"
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
type UserInfo struct {
	Nick string `json:"nickname"`
	Snap string `json:"snap"`
}
type GroupInfo struct {
	Groupid string     `json:"groupid"`
	Creator string     `json:"creator"`
	Name    string     `json:"name"`
	Notice  string     `json:"notice"`
	Snap    string     `json:"snap"`
	Members []UserInfo `json:"members"`
}
type JsonResponseData struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Data    GroupInfo `json:"data"`
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
	name := req.FormValue("name")

	if len(name) < 1 || len(creator) < 2 || len(members) < 2 {
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
	snapurl := snap.GenGroupSnap(bmembers, strID)
	userinfo.UserInfoCh <- bmembers

	client.HMSet(key, "creator", creator, "name", name, "notice", "", "snap", snapurl)
	for _, v := range users {
		client.SAdd(mkey, v)
		skey := "mygroups_" + v
		client.SAdd(skey, strID)
	}

	client.Incr("groupID")
	client.Close()
	//groupinfo := GroupInfo{strID, creator, name, "", snapurl, users}
	var wrong []UserInfo
	groupinfo := GroupInfo{strID, creator, name, "", snapurl, wrong}
	jsonres := JsonResponseData{1, "OK", groupinfo}
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
	userinfo.UserInfoCh <- bmembers

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
	howmany, _ := client.SMembers(mkey)
	if len(howmany) <= 9 {
		howmanys := strings.Join(howmany, ",")
		snapurl := snap.GenGroupSnap(howmanys, groupid)
		client.HSet("group_"+groupid, "snap", snapurl)
	}

	client.Close()

	chat.Mutex.Lock()
	delete(chat.Groups, groupid)
	chat.Mutex.Unlock()

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
	howmany, _ := client.SMembers(mkey)
	if len(howmany) < 9 {
		howmanys := strings.Join(howmany, ",")
		snapurl := snap.GenGroupSnap(howmanys, groupid)
		client.HSet("group_"+groupid, "snap", snapurl)
	}

	client.Close()

	chat.Mutex.Lock()
	delete(chat.Groups, groupid)
	chat.Mutex.Unlock()

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
	var agroups []GroupInfo
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
	for _, v := range groups {
		kkey := "group_" + v
		ls, _ := client.HMGet(kkey, "creator", "name", "notice", "snap")
		var temp GroupInfo
		temp.Groupid = v
		for k, p := range ls {
			switch {
			case k == 0:
				temp.Creator = p
			case k == 1:
				temp.Name = p
			case k == 2:
				temp.Notice = p
			case k == 3:
				temp.Snap = p
			}
		}
		gkey := "groupmembers_" + v
		members, _ := client.SMembers(gkey)
		for _, member := range members {
			lss, _ := client.HMGet("user_"+member, "nick", "snap")
			var tempp UserInfo
			for k, p := range lss {
				switch {
				case k == 0:
					tempp.Nick = p
				case k == 1:
					tempp.Snap = p
				}
			}
			temp.Members = append(temp.Members, tempp)

		}
		agroups = append(agroups, temp)
	}
	client.Close()
	type MyResopnse struct {
		JsonResponse
		Groups []GroupInfo `json:"data"`
	}
	jsonres := MyResopnse{JsonResponse{1, "OK"}, agroups}
	b, _ := json.Marshal(jsonres)
	io.WriteString(w, string(b))
	return

}
func main() {

	runtime.GOMAXPROCS(4)
	log.SetFlags(log.Lshortfile)
	chat.InitGroup()
	chat.InitRedis()
	userinfo.InitUserCh()

	// websocket server
	server := chat.NewServer("/entry")
	go server.Listen()
	go chat.RecGrpMsgTrd(server)
	//	go chat.HeaartbeatTrd(server)
	go userinfo.GetUserinfo()

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
