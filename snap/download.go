package main

import "os"
import "io"

//import "fmt"
import "strings"
import "net/http"

const PATH = "/live/www/html/emovideo/"

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

//func main() {
//	var urls = [...]string{"https://imgcdn.66boss.com/imagesu/avatar/1470558849_10000016530.jpg",
//		"https://imgcdn.66boss.com/imagesu/avatar_temp/default.jpg"}
//	outputs := Download(urls[:])
//	fmt.Println(outputs)
//}
