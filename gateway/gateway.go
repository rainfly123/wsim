package gateway

import (
	"log"
)

func TransmitMessage(touser string, msg string) bool {
	server := searchServer()
	if len(server) <= 0 {
		return false
	}
	var temp = Transmit{server, touser, msg}
	receiveChan <- temp
	return true
}

func searchServer(touser string) string {
	var client *redis.Client
	var ok bool
	var has bool

	client, ok = Clients.Get()
	if ok != true {
		log.Panic("redis error")
		return
	}
	servers, _ := client.SMembers("servers")
	for _, v := range servers {
		has, _ = client.SIsMember(v+"_users", touser)
		if has == true {
			return v
		}
	}
	defer client.Close()
	return ""
}

func ReportOnline() {
	var client *redis.Client
	var ok bool

	client, ok = Clients.Get()
	if ok != true {
		log.Panic("redis error")
		return
	}
	client.SAdd("servers", LocalIPRecvAddr)
	client.Close()
}

func ReportOffline() {
	var client *redis.Client
	var ok bool

	client, ok = Clients.Get()
	if ok != true {
		log.Panic("redis error")
		return
	}
	client.SRem("servers", LocalIPRecvAddr)
	client.Close()
}

func ReportUserOnline(touser string) {
	var client *redis.Client
	var ok bool

	client, ok = Clients.Get()
	if ok != true {
		log.Panic("redis error")
		return
	}
	client.SAdd(LocalIPRecvAddr+"_users", touser)
	client.Close()

}

func ReportUserOffline(touser string) {
	var client *redis.Client
	var ok bool

	client, ok = Clients.Get()
	if ok != true {
		log.Panic("redis error")
		return
	}
	client.SRem(LocalIPRecvAddr+"_users", touser)
	client.Close()
}
