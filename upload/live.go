package main

import (
	"fmt"
	//"io/ioutil"
	//"net/http"
	//	"os"
	"os/exec"
	//"path"
	"strings"
	//"time"
)

func Checkvideo(origin string) {

	index := strings.LastIndex(origin, ".")
	if index > 0 {
		dest := origin[0:index+1] + "jpg"
		//var args = []string{"-i", origin, "-vframes", "1", "-vf", "crop=iw:iw*9/16", "-f", "image2", "-y", dest}
		var args = []string{"-i", origin, "-vframes", "1", "-f", "image2", "-y", dest}
		cmd := exec.Command("ffmpeg", args[0:]...)
		//cmd.Stdout = os.Stdout
		//cmd.Stderr = os.Stderr
		err := cmd.Run()
		logger.Println(err)
	}
}

func Check_Thread() {
	for video := range Channel {
		go Checkvideo(video)
	}
}

var Channel chan string
