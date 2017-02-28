package userinfo

import "menteslibres.net/gosexy/redis"
import "log"

//import "fmt"
import "strings"
import "encoding/json"
import "io/ioutil"
import "net/http"

var UserInfoCh chan string
var client *redis.Client
var host = "127.0.0.1"
var port = uint(6379)

const SNAPURL = "https://api.66boss.com/ucenter/userinfo/avatar?user_id="
const NICK = "https://api.66boss.com/ucenter/userinfo/nick?user_id="

func InitUserCh() {
	UserInfoCh = make(chan string)
	client = redis.New()
	err := client.Connect(host, port)
	client.Command(nil, "SELECT", 1)
	log.Println("userinfo inited", err)
}

func GetUserinfo() {
	for users := range UserInfoCh {

		res, err := http.Get(SNAPURL + users)
		if err != nil {
			log.Println(err)
			continue
		}
		detail, err := ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Println(err)
			continue
		}
		all := make(map[string]interface{})
		json.Unmarshal(detail, &all)
		result := all["result"]
		snaps := result.(map[string]interface{})
		all_users := strings.Split(users, ",")
		for _, user := range all_users {
			url := snaps[user].(string)
			client.HSet("user_"+user, "snap", url)
		}
		res, err = http.Get(NICK + users)
		if err != nil {
			log.Println(err)
			continue
		}
		detail, err = ioutil.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Println(err)
			continue
		}
		//all := make(map[string]interface{})
		json.Unmarshal(detail, &all)
		result = all["result"]
		snaps = result.(map[string]interface{})
		all_users = strings.Split(users, ",")
		for _, user := range all_users {
			nick := snaps[user].(string)
			client.HSet("user_"+user, "nick", nick)
		}

	}
}

/*
func main() {
	InitUserCh()
	go GetUserinfo()
	UserInfoCh <- "100000035,100000034"
	for {
		//		fmt.Println(client.HMGet("100000034", "nick", "snap"))
	}

}
*/
