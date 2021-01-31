package panel

import (
	"errors"
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
)

const nameSuffix string = "_marked"

func readImg(f fyne.URIReadCloser, globalData *Data) (*canvas.Image, string, error) {
	if f == nil {
		log.Println("Cancelled")
		return nil, "", errors.New("file not found")
	}
	img := loadImage(f)
	path := f.URI().String()
	globalData.inputName = f.URI().Name()
	globalData.inputExt = f.URI().Extension()
	err := f.Close()
	if err != nil {
		fyne.LogError("Failed to close stream", err)
		return nil, "", errors.New("file not found")
	}
	if img == nil {
		return nil, "", errors.New("file not found")
	}
	img.FillMode = canvas.ImageFillContain
	return img, path, nil

}
func loadImage(f fyne.URIReadCloser) *canvas.Image {
	data, err := ioutil.ReadAll(f)
	if err != nil {
		fyne.LogError("Failed to load image data", err)
		return nil
	}
	res := fyne.NewStaticResource(f.URI().Name(), data)
	return canvas.NewImageFromResource(res)
}
func writeImg(data []byte, path string, topWindow fyne.Window, globalData *Data, ch chan float64) error {
	defer func(ch chan float64) {
		ch <- 1
		close(ch)
	}(ch)
	file, _ := pathExists(path)
	skip := false
	if file {
		dialog.NewConfirm("confirm ?", "file exist ,do you want overrides it ?", func(b bool) {
			if !b {
				skip = true
			}
		}, topWindow)
	}
	if !skip {
		idx := len(globalData.inputExt)
		finalPath := fmt.Sprintf("%s%c%s%s%s", path[7:], os.PathSeparator, globalData.inputName[0:len(globalData.inputName)-idx], nameSuffix, globalData.inputExt)
		fmt.Println(
			finalPath,
		)
		err := ioutil.WriteFile(finalPath, data, 0755)
		return err
	}
	return errors.New("没有生成图片")
}
func rgbGradient(x, y, w, h int) color.Color {
	g := int(float32(x) / float32(w) * float32(255))
	b := int(float32(y) / float32(h) * float32(255))
	return color.NRGBA{uint8(255 - b), uint8(g), uint8(b), 0xff}
}

/*
   判断文件或文件夹是否存在
   如果返回的错误为nil,说明文件或文件夹存在
   如果返回的错误类型使用os.IsNotExist()判断为true,说明文件或文件夹不存在
   如果返回的错误为其它类型,则不确定是否在存在
*/
func pathExists(path string) (bool, error) {

	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
