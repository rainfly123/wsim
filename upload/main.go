package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	//"sort"
	"strconv"
	"strings"
	//	"time"
)

var logger *log.Logger

type JsonResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Url     string `json:"data"`
}

func writev1Handle(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		io.WriteString(w, fmt.Sprintf("<html><head><title>我的第一个页面</title></head><body><form action=\"writev1\" method=\"post\" enctype=\"multipart/form-data\"><label>上传音频</label><input type=\"file\" name='file'/><br/><label><input type=\"submit\" value=\"上传音频\"/></label></form></body></html>"))
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
		mp3file := Checkaudio(uuidFile)
		var newFile string
		newFile = mp3file
		duration := GetDuration(mp3file)
		index := strings.LastIndex(mp3file, ".")
		if index > 0 {
			newFile = mp3file[0:index] + "-" + duration + ".mp3"
			os.Rename(mp3file, newFile)
		}
		jsonres := JsonResponse{0, "OK", (ACCESS_URL + path.Base(newFile))}
		b, _ := json.Marshal(jsonres)
		io.WriteString(w, string(b))
		//Channel <- uuidFile
	}
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
		withwh := Resize(uuidFile)
		jsonres := JsonResponse{0, "OK", (ACCESS_URL + withwh)}
		b, _ := json.Marshal(jsonres)
		io.WriteString(w, string(b))
		//Channel <- uuidFile
	}
}

type Meta struct {
	Url  string `json:"url"`
	Snap string `json:"snap"`
}
type mJsonResponse struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Metainfo Meta   `json:"data"`
}

func writev3Handle(w http.ResponseWriter, req *http.Request) {
	var uuidFile string
	if req.Method == "GET" {
		io.WriteString(w, fmt.Sprintf("<html><head><title>我的第一个页面</title></head><body><form action=\"writev3\" method=\"post\" enctype=\"multipart/form-data\"><label>上传视频</label><input type=\"file\" name='file'/><br/><label><input type=\"submit\" value=\"上传视频\"/></label></form></body></html>"))
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
		//Channel <- uuidFile
		snapfile := Checkvideo(uuidFile)
		mygod := WidthHeightFile(snapfile)
		index := strings.LastIndex(mygod, ".")
		var newFile string
		if index > 0 {
			newFile = mygod[0:index] + path.Ext(uuidFile)
			os.Rename(uuidFile, (UPLOAD_VIDEO_PATH + newFile))
		}
		jsonres := JsonResponse{0, "Succeeded", (ACCESS_VIDEO_URL + newFile)}
		b, _ := json.Marshal(jsonres)
		io.WriteString(w, string(b))
	}

}

func main() {

	logfile, _ := os.OpenFile("/var/log/upload.log", os.O_RDWR|os.O_CREATE, 0)
	logger = log.New(logfile, "\n", log.Ldate|log.Ltime|log.Lshortfile)
	//Channel = make(chan string, 100)
	//go Check_Thread()

	http.HandleFunc("/writev1", writev1Handle)
	http.HandleFunc("/writev2", writev2Handle)
	http.HandleFunc("/writev3", writev3Handle)

	//	http.Handle("/", http.FileServer(http.Dir("/root/git/weibo/upload")))

	if err := http.ListenAndServe(":6070", nil); err != nil {
	}
}
