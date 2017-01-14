package main

import "os"
import "io"

import "fmt"
import "strings"
import "encoding/json"
import "net/http"
import "io/ioutil"

const PATH = "/live/www/html/emovideo/"
const SNAPURL = "http://www.66boss.com/app/user.php?act=user_avatar&user_id="

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func Download(urls []string) []string {
	var outputs []string
	for _, url := range urls {
		if strings.Contains(url, "default") {
			outputs = append(outputs, (PATH + "default.jpg"))
			continue
		}

		index := strings.LastIndex(url, "_") + 1
		uuid := PATH + url[index:]
		if Exist(uuid) {
			outputs = append(outputs, uuid)
			continue
		}

		resp, err := http.Get(url)
		if err != nil {
			// handle error
			return outputs
		}
		fW, err := os.Create(uuid)
		if err != nil {
			return outputs
		}
		_, err = io.Copy(fW, resp.Body)
		fW.Close()
		resp.Body.Close()
		outputs = append(outputs, uuid)
	}
	return outputs
}

func GetURLs(users string) []string {
	var urls []string
	res, err := http.Get(SNAPURL + users)
	if err != nil {
		return urls
	}
	detail, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return urls
	}
	var all interface{}
	eb := json.Unmarshal(detail, &all)
	fmt.Println(all, detail, eb)
	data := all.(map[string]interface{})
	all_users := strings.Split(users, ",")
	for _, user := range all_users {
		url := data[user].(string)
		fmt.Println(url)
		urls = append(urls, url)
	}
	return urls
}

func main() {
	//	var urls = [...]string{"https://imgcdn.66boss.com/imagesu/avatar/1470558849_10000016530.jpg",
	//		"https://imgcdn.66boss.com/imagesu/avatar_temp/default.jpg"}
	//	outputs := Download(urls[:])
	//	fmt.Println(outputs)
	var cc = "1,1000001653"
	fmt.Println(GetURLs(cc))
}
