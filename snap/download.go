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
	for i, ch := range detail {

		switch {
		case ch > '~':
			detail[i] = ' '
		case ch == '\r':
		case ch == '\n':
		case ch == '\t':
		case ch < ' ':
			detail[i] = ' '
		}
	}
	res.Body.Close()
	if err != nil {
		return urls
	}
	all := make(map[string]interface{})
	json.Unmarshal(detail, &all)
	all_users := strings.Split(users, ",")
	for _, user := range all_users {
		url := all[user].(string)
		urls = append(urls, url)
	}
	return urls
}

/*
func main() {
	//	var urls = [...]string{"https://imgcdn.66boss.com/imagesu/avatar/1470558849_10000016530.jpg",
	//		"https://imgcdn.66boss.com/imagesu/avatar_temp/default.jpg"}
	//	outputs := Download(urls[:])
	//	fmt.Println(outputs)
	var cc = "1,1000001653"
	d := GetURLs(cc)
	fmt.Println(Download(d))
}
*/
