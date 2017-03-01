package main

import (
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path"
	"strings"
)

func Resize(name string) string {
	decoder := jpeg.Decode
	var JPG bool

	if !strings.HasSuffix(name, ".jpg") && !strings.HasSuffix(name, ".png") {
		return path.Base(name)
	}

	JPG = strings.HasSuffix(name, ".jpg")
	if !JPG {
		decoder = png.Decode
	}

	file, err := os.Open(name)
	if err != nil {
		return path.Base(name)
	}

	img, err := decoder(file)
	if err != nil {
		return path.Base(name)
	}
	file.Close()
	var m image.Image
	fileinfo, _ := os.Stat(name)
	if fileinfo.Size()/1024 > 1024 {
		m = resize.Resize(600, 0, img, resize.Lanczos3)
	} else {
		m = img
	}

	ext := path.Ext(name)
	bounds := m.Bounds()
	min, max := bounds.Min, bounds.Max
	height, width := max.Y-min.Y, max.X-min.X

	index := strings.LastIndex(name, ".")
	var temp string = fmt.Sprintf("%s-%dx%d%s", name[:index], width, height, ext)
	out, err := os.Create(temp)
	if err != nil {
		return path.Base(name)
	}
	// write new image to file
	if JPG {
		jpeg.Encode(out, m, nil)
	} else {
		png.Encode(out, m)
	}
	out.Close()
	//os.Rename(temp, name)
	return path.Base(temp)
}

/*
func main() {
	fmt.Println(os.Args[1])
	fmt.Println(Resize(os.Args[1]))
}
*/
