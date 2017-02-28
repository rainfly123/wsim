package chat

import "sync"
import "menteslibres.net/gosexy/redis"
import "log"
import "time"

var Groups map[string][]string
var Mutex sync.Mutex
var GroupMsgCh chan *InPut

const MEMBERS = "http://www.66boss.com/app/tribe.php?act=tribe_user_id&tribe_id="

func InitGroup() {
	Groups = make(map[string][]string)
	GroupMsgCh = make(chan *InPut)
	log.Println("Groups Inited", Groups)
}

func checkGroup(input *InPut, server *Server) {

	var client *redis.Client
	var ok bool

	groupid := input.Touserid
	key := "groupmembers_" + groupid

	client, ok = Clients.Get()
	if ok != true {
		log.Panic("redis error")
		return
	}

	users, _ := client.SMembers(key)
	defer client.Close()

	Mutex.Lock()
	Groups[groupid] = users
	Mutex.Unlock()

	for _, user := range users {
		if user == input.Fromuserid {
			continue
		}
		lockUsers.RLock()
		client, online := server.users[user]
		lockUsers.RUnlock()
		output := NewOutput(input)
		if online {
			//output := NewOutput(input)
			client.Write([]byte(output.String()))
		} else {
			//user is offline
			//log.Println("user is offline", user)
			PushOfflineMsg(user, output.String())
		}
	}
}

func RecGrpMsgTrd(server *Server) {
	for input := range GroupMsgCh {
		groupid := input.Touserid

		Mutex.Lock()
		users, ok := Groups[groupid]
		Mutex.Unlock()

		if ok {
			for _, user := range users {
				if user == input.Fromuserid {
					continue
				}
				lockUsers.RLock()
				client, online := server.users[user]
				lockUsers.RUnlock()
				output := NewOutput(input)
				if online {
					//output := NewOutput(input)
					client.Write([]byte(output.String()))
				} else {
					//user is offline
					//log.Println("user is offline", user)
					PushOfflineMsg(user, output.String())
				}
			}
		} else {
			// no this group,  checking ...
			log.Println("group not exists check...")
			go checkGroup(input, server)
		}
	}
}

func HeaartbeatTrd(server *Server) {
	for {
		var users map[string]*Client
		time.Sleep(30 * time.Second)

		lockUsers.RLock()
		users = server.users
		lockUsers.RUnlock()

		for _, client := range users {
			client.Write([]byte("heartbeat_argument"))
			//log.Println(client.userid)
		}
	}
}
