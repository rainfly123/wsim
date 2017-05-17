package gateway

import (
	"log"
)

func SearchServer(touser string) string {
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
	client.Close()
}

func ReportOnline() {
	var client *redis.Client
	var ok bool

	client, ok = Clients.Get()
	if ok != true {
		log.Panic("redis error")
		return
	}
	client.SAdd("servers", LocalIP)
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
	client.SRem("servers", LocalIP)
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
	client.SAdd(LocalIP+"_users", touser)
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
	client.SRem(LocalIP+"_users", touser)
	client.Close()
}
