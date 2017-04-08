package snap

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

var POSITION = [10][9]image.Rectangle{
	{},
	{},
	{},
	{0: image.Rect(95, 10, 245, 160), 1: image.Rect(10, 170, 160, 320), 2: image.Rect(170, 170, 320, 320)},
	{0: image.Rect(10, 10, 160, 160), 1: image.Rect(170, 10, 330, 160), 2: image.Rect(10, 170, 170, 330), 3: image.Rect(170, 170, 330, 330)},
	{0: image.Rect(65, 65, 165, 165), 1: image.Rect(175, 65, 275, 165), 2: image.Rect(10, 175, 110, 275), 3: image.Rect(120, 175, 220, 275), 4: image.Rect(230, 175, 330, 275)},
	{0: image.Rect(10, 65, 110, 165), 1: image.Rect(120, 65, 220, 165), 2: image.Rect(230, 65, 330, 165), 3: image.Rect(10, 175, 110, 275), 4: image.Rect(120, 175, 220, 275), 5: image.Rect(230, 175, 330, 275)},
	{0: image.Rect(10, 10, 110, 110), 1: image.Rect(120, 10, 220, 210), 2: image.Rect(230, 10, 330, 110), 3: image.Rect(10, 120, 110, 220), 4: image.Rect(120, 120, 220, 220), 5: image.Rect(230, 120, 330, 220), 6: image.Rect(10, 230, 110, 330)},
	{0: image.Rect(10, 10, 110, 110), 1: image.Rect(120, 10, 220, 210), 2: image.Rect(230, 10, 330, 110), 3: image.Rect(10, 120, 110, 220), 4: image.Rect(120, 120, 220, 220), 5: image.Rect(230, 120, 330, 220), 6: image.Rect(10, 230, 110, 330), 7: image.Rect(120, 230, 220, 330)},
	{0: image.Rect(10, 10, 110, 110), 1: image.Rect(120, 10, 220, 210), 2: image.Rect(230, 10, 330, 110), 3: image.Rect(10, 120, 110, 220), 4: image.Rect(120, 120, 220, 220), 5: image.Rect(230, 120, 330, 220), 6: image.Rect(10, 230, 110, 330), 7: image.Rect(120, 230, 220, 330), 8: image.Rect(230, 230, 330, 330)},
}

func getUUID() string {
	f, _ := os.OpenFile("/dev/urandom", os.O_RDONLY, 0)
	b := make([]byte, 8)
	f.Read(b)
	f.Close()
	uuid := fmt.Sprintf("_%x%x%x%x.jpg", b[0:1], b[2:3], b[4:5], b[6:7])
	return uuid
}

func Resize(pictures []string) ([]image.Image, int) {
	var outputs []image.Image
	total := len(pictures)
	var size uint = 100
	if total <= 4 {
		size = 150
	}
	for _, picture := range pictures {
		var PNG bool
		var img image.Image

		PNG = strings.HasSuffix(picture, ".png")
		file, err := os.Open(picture)
		if err != nil {
			return outputs, len(outputs)
		}
		if PNG {
			img, err = png.Decode(file)
		} else {
			img, err = jpeg.Decode(file)
		}
		if err != nil {
			file.Close()
			file, err = os.Open(picture)
			if !PNG {
				img, err = png.Decode(file)
			} else {
				img, err = jpeg.Decode(file)
			}
		}
		file.Close()

		m := resize.Resize(size, size, img, resize.Lanczos3)
		outputs = append(outputs, m)
	}
	return outputs, len(outputs)
}

func GenGroupSnap(users string, group string) string {
	filename := group + getUUID()
	dstfile := PATH + filename
	file, err := os.Create(dstfile)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	urls := GetURLs(users)
	files := Download(urls)

	m := image.NewRGBA(image.Rect(0, 0, 340, 340))
	blue := color.RGBA{200, 200, 200, 255}
	draw.Draw(m, m.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)

	outputs, _ := Resize(files)
	fmt.Println(outputs)
	total := len(outputs)
        if total >= 10 {
            total = 9
        }
	for i, output := range outputs {
                if i >= 9 {
                    break
                }
		tmp := POSITION[total][i]
		draw.Draw(m, tmp, output, output.Bounds().Min, draw.Over)
	}
	jpeg.Encode(file, m, nil)
	return VPATH + filename
}

/*
func main() {
	fmt.Println(GenGroupSnap("1000006331,1000006123,1000006340", "333"))
	//	fmt.Println(GenGroupSnap("1000001653,1000001653,1000001653", "333"))
}
*/
