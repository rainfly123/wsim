package main

import (
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

func Resize(pictures []string) ([]image.Image, int) {
	var outputs []image.Image

	for _, picture := range pictures {
		decoder := jpeg.Decode
		var JPG bool
		JPG = strings.HasSuffix(picture, ".jpg")
		if !JPG {
			decoder = png.Decode
		}

		file, err := os.Open(picture)
		if err != nil {
			return outputs, len(outputs)
		}

		img, err := decoder(file)
		if err != nil {
			return outputs, len(outputs)
		}
		file.Close()

		m := resize.Resize(100, 100, img, resize.Lanczos3)
		outputs = append(outputs, m)
	}
	return outputs, len(outputs)
}

func main() {
	file, err := os.Create("dst.jpg")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	//	file1, err := os.Open("20.jpg")
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	defer file1.Close()
	//	img, _ := jpeg.Decode(file1)

	//	jpg := image.NewGray(image.Rect(0, 0, 610, 610)) //NewGray
	//	for x := 0; x < 610; x++ {
	//		for y := 0; y < 610; y++ {
	//			jpg.Set(x, y, color.Gray{200}) //设定alpha图片的透明度
	//		}
	//	}

	m := image.NewRGBA(image.Rect(0, 0, 340, 340))
	blue := color.RGBA{200, 200, 200, 255}
	draw.Draw(m, m.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)

	//jpg := image.NewGray(img.Bounds())                                 //NewGray

	//	draw.Draw(jpg, jpg.Bounds(), img, img.Bounds().Min, draw.Src) //原始图片转换为灰色图片

	var files = []string{"/live/www/html/emovideo/10000016530.jpg", "/live/www/html/emovideo/default.jpg"}
	outputs, _ := Resize(files)
	for i, output := range outputs {
		if i == 0 {
			draw.Draw(m, m.Bounds(), output, output.Bounds().Min, draw.Over)
		} else {
			draw.Draw(m, image.Rect(100, 100, 200, 200), output, output.Bounds().Min, draw.Over)
		}
	}
	draw.Draw(m, image.Rect(200, 200, 300, 300), outputs[0], outputs[0].Bounds().Min, draw.Over)
	jpeg.Encode(file, m, nil)

}
