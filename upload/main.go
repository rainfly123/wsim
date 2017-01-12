package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	//"sort"
	"strconv"
	//	"strings"
	//	"time"
)

var logger *log.Logger

type JsonResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Url     string `json:"data"`
}

func writev2Handle(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		io.WriteString(w, fmt.Sprintf("<html><head><title>我的第一个页面</title></head><body><form action=\"writev2\" method=\"post\" enctype=\"multipart/form-data\"><label>上传图片</label><input type=\"file\" name='file'/><br/><label><input type=\"submit\" value=\"上传图片\"/></label></form></body></html>"))
	} else {
		file, head, err := req.FormFile("file")
		if err != nil {
			jsonres := JsonResponse{1, "参数错误", ""}
			b, _ := json.Marshal(jsonres)
			io.WriteString(w, string(b))
			return
		}
		defer file.Close()
		temp := getFileName(head.Filename)
		uuidFile := UPLOAD_PATH + temp
		fW, err := os.Create(uuidFile)
		if err != nil {
			jsonres := JsonResponse{2, "系统错误", ""}
			b, _ := json.Marshal(jsonres)
			io.WriteString(w, string(b))
			return
		}
		_, err = io.Copy(fW, file)
		if err != nil {
			jsonres := JsonResponse{2, "系统错误", ""}
			b, _ := json.Marshal(jsonres)
			io.WriteString(w, string(b))
			return
		}
		Resize(uuidFile)
		jsonres := JsonResponse{0, "OK", (ACCESS_URL + temp)}
		b, _ := json.Marshal(jsonres)
		io.WriteString(w, string(b))
		Channel <- uuidFile
	}
}

func writev3Handle(w http.ResponseWriter, req *http.Request) {
	var uuidFile string
	if req.Method == "GET" {
		io.WriteString(w, fmt.Sprintf("<html><head><title>我的第一个页面</title></head><body><form action=\"writev3\" method=\"post\" enctype=\"multipart/form-data\"><label>上传图片</label><input type=\"file\" name='file'/><br/><label><input type=\"submit\" value=\"上传视频\"/></label></form></body></html>"))
	} else {
		filesize, _ := strconv.Atoi(req.Header.Get("Content-Length"))
		if filesize/1024/1024 > 20 {
			jsonres := JsonResponse{3, "视频文件过大", ""}
			b, _ := json.Marshal(jsonres)
			io.WriteString(w, string(b))
			return
		}
		file, head, err := req.FormFile("file")
		if err != nil {
			jsonres := JsonResponse{1, "参数错误", ""}
			b, _ := json.Marshal(jsonres)
			io.WriteString(w, string(b))
			return
		}
		defer file.Close()

		temp := getFileName(head.Filename)
		uuidFile = UPLOAD_VIDEO_PATH + temp
		fW, err := os.Create(uuidFile)
		if err != nil {
			jsonres := JsonResponse{2, "系统错误", ""}
			b, _ := json.Marshal(jsonres)
			io.WriteString(w, string(b))
			return
		}
		_, err = io.Copy(fW, file)
		if err != nil {
			jsonres := JsonResponse{2, "系统错误", ""}
			b, _ := json.Marshal(jsonres)
			io.WriteString(w, string(b))
			return
		}
		jsonres := JsonResponse{0, "Succeeded", (ACCESS_VIDEO_URL + temp)}
		b, _ := json.Marshal(jsonres)
		io.WriteString(w, string(b))
		Channel <- uuidFile
	}

}

func main() {

	logfile, _ := os.OpenFile("/var/log/upload.log", os.O_RDWR|os.O_CREATE, 0)
	logger = log.New(logfile, "\n", log.Ldate|log.Ltime|log.Lshortfile)
	Channel = make(chan string, 100)
	go Check_Thread()

	http.HandleFunc("/writev2", writev2Handle)
	http.HandleFunc("/writev3", writev3Handle)

	//	http.Handle("/", http.FileServer(http.Dir("/root/git/weibo/upload")))

	if err := http.ListenAndServe(":9090", nil); err != nil {
	}
}
