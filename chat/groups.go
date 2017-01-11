package chat

import "sync"
import "io/ioutil"
import "net/http"
import "log"
import "time"
import "strings"

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

	groupid := input.Touserid
	url := MEMBERS + groupid //5874
	res, err := http.Get(url)
	if err != nil {
		log.Println("Error Can't connect to www.66boss.com")
		return
	}
	detail, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if len(detail) <= 2 {
		log.Println("Error Read www.66boss.com")
		return
	}
	result := string(detail)
	result = strings.Replace(result, "[", "", -1)
	result = strings.Replace(result, "]", "", -1)
	users := strings.Split(result, ",")

	Mutex.Lock()
	Groups[groupid] = users
	Mutex.Unlock()

	for _, user := range users {
		lockUsers.RLock()
		client, online := server.users[user]
		lockUsers.RUnlock()
		output := NewOutput(input)
		if online {
			//output := NewOutput(input)
			client.Write(output.Bytes())
		} else {
			//user is offline
			//log.Println("user is offline", user)
			PushOfflineMsg(input.Touserid, output.String())
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
				lockUsers.RLock()
				client, online := server.users[user]
				lockUsers.RUnlock()
				output := NewOutput(input)
				if online {
					//output := NewOutput(input)
					client.Write(output.Bytes())
				} else {
					//user is offline
					//log.Println("user is offline", user)
					PushOfflineMsg(input.Touserid, output.String())
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
		time.Sleep(10 * time.Second)

		lockUsers.RLock()
		users = server.users
		lockUsers.RUnlock()

		for _, client := range users {
			client.Write([]byte("heartbeat_argument"))
			log.Println(client.userid)
		}
	}
}
