package chat

import "sync"
import "log"

var Groups map[string][]string
var Mutex sync.Mutex
var GroupMsgCh chan *InPut

func InitGroup() {
	Groups := make(map[string][]string)
	log.Println("Groups Inited", Groups)
}

func checkGroup(input *InPut) {

}

func RecGrpMsgTrd(server *Server) {
	for input := range GroupMsgCh {
		groupid := input.Touserid

		Mutex.Lock()
		users, ok := Groups[groupid]
		Mutex.Unlock()

		if ok {
			for _, user := range users {
				client, online := server.users[user]
				if online {
					output := NewOutput(input)
					client.Write(output.Bytes())
				} else {
					//user is offline
				}
			}
		} else {
			// no this group,  checking ...
			go checkGroup(input)
		}
	}
}
