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
	"github.com/tcolgate/mp3"
	"image"
	"image/jpeg"
	"os"
	"path"
	"strconv"
	"time"
)

func GetDuration(file string) string {
	r, err := os.Open(file)
	if err != nil {
		return "0"
	}

	d := mp3.NewDecoder(r)
	var f mp3.Frame
	var du time.Duration
	for {

		if err := d.Decode(&f); err != nil {
			break
		}
		du += f.Duration()

	}
	r.Close()
	return strconv.Itoa((int(du.Seconds() + 0.5)))
}

func WidthHeightFile(picture string) string {
	total := len(picture)
	if total < 5 {
		return ""
	}
	var img image.Image
	file, err := os.Open(picture)
	if err != nil {
		return ""
	}

	img, err = jpeg.Decode(file)
	file.Close()

	bounds := img.Bounds()
	min, max := bounds.Min, bounds.Max
	height, width := max.Y-min.Y, max.X-min.X

	index := strings.LastIndex(picture, ".")
	ext := path.Ext(picture)
	var temp string = fmt.Sprintf("%s-%dx%d%s", picture[:index], width, height, ext)
	out, err := os.Create(temp)
	if err != nil {
		return path.Base(picture)
	}
	// write new image to file
	jpeg.Encode(out, img, nil)
	out.Close()
	//os.Rename(temp, name)
	return path.Base(temp)
}

func Checkvideo(origin string) string {

	index := strings.LastIndex(origin, ".")
	if index > 0 {
		dest := origin[0:index+1] + "jpg"
		//var args = []string{"-i", origin, "-vframes", "1", "-vf", "crop=iw:iw*9/16", "-f", "image2", "-y", dest}
		var args = []string{"-i", origin, "-vframes", "1", "-f", "image2", "-y", dest}
		cmd := exec.Command("ffmpeg", args[0:]...)
		//cmd.Stdout = os.Stdout
		//cmd.Stderr = os.Stderr
		err := cmd.Run()
		fmt.Println(err)
		return dest
	}
	return ""
}
func Checkaudio(origin string) string {

	index := strings.LastIndex(origin, ".")
	if index > 0 {
		dest := origin[0:index+1] + "mp3"
		//var args = []string{"-i", origin, "-vframes", "1", "-vf", "crop=iw:iw*9/16", "-f", "image2", "-y", dest}
		var args = []string{"-i", origin, "-y", dest}
		cmd := exec.Command("ffmpeg", args[0:]...)
		//cmd.Stdout = os.Stdout
		//cmd.Stderr = os.Stderr
		err := cmd.Run()
		fmt.Println(err)
		return dest
	}
	return ""
}

//func Check_Thread() {
//	for video := range Channel {
//		go Checkvideo(video)
//	}
//}

//var Channel chan string

/*func main() {
	fmt.Println(WidthHeightFile(os.Args[1]))
}*/
