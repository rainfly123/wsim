package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	UPLOAD_PATH string = "/live/www/html/weibo_pic/"
	//ACCESS_URL  string = "https://weibocdn.66boss.com/weibo_pic/"
	ACCESS_URL        string = "http://livecdn.66boss.com/weibo_pic/"
	UPLOAD_VIDEO_PATH string = "/live/www/html/weibo_video/"
	ACCESS_VIDEO_URL  string = "http://livecdn.66boss.com/weibo_video/"
	//ACCESS_VIDEO_URL  string = "https://weibocdn.66boss.com/weibo_video/"

)

func getUUID() string {
	f, _ := os.OpenFile("/dev/urandom", os.O_RDONLY, 0)
	b := make([]byte, 16)
	f.Read(b)
	f.Close()
	uuid := fmt.Sprintf("%x%x%x%x%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return uuid
}

func getFileName(name string) string {

	var temp string = "error"
	i := strings.LastIndex(name, ".")
	if i > 0 {
		uuid := getUUID()
		temp = uuid + name[i:]
	}
	return temp
}
