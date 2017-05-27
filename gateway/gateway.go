package gateway

import (
	"log"
	//	"menteslibres.net/gosexy/redis"
)

func TransmitMessage(touser string, msg string) bool {
	server := searchServer(touser)
	if len(server) <= 0 {
		return false
	}
	var temp = Transmit{server, touser, msg}
	receiveChan <- temp
	return true
}

func searchServer(touser string) string {
	var has bool

	client, ok := clients.Get()
	if ok != true {
		log.Panic("redis error")
		return ""
	}
	servers, _ := client.SMembers("servers")
	for _, v := range servers {
		if v == LocalIPRecvAddr {
			continue
		}
		has, _ = client.SIsMember(v+"_users", touser)
		if has == true {
			return v
		}
	}
	defer client.Close()
	return ""
}

func ReportOnline() {

	client, ok := clients.Get()
	if ok != true {
		log.Panic("redis error")
		return
	}
	client.SAdd("servers", LocalIPRecvAddr)
	client.Close()
}

func ReportOffline() {

	client, ok := clients.Get()
	if ok != true {
		log.Panic("redis error")
		return
	}
	client.SRem("servers", LocalIPRecvAddr)
	client.Close()
}

func ReportUserOnline(touser string) {

	client, ok := clients.Get()
	if ok != true {
		log.Panic("redis error")
		return
	}
	client.SAdd(LocalIPRecvAddr+"_users", touser)
	client.Close()

}

func ReportUserOffline(touser string) {

	client, ok := clients.Get()
	if ok != true {
		log.Panic("redis error")
		return
	}
	client.SRem(LocalIPRecvAddr+"_users", touser)
	client.Close()
}
