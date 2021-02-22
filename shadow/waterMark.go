package shadow

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"strings"
	"unicode"

	"github.com/goki/freetype"
	"github.com/goki/freetype/truetype"
)

const minHeight int = 20

//DealImg 处理图片，加水印
func DealImg(data []byte, watermark string, ext string, ch chan float64) ([]byte, error) {
	defer func(ch chan float64) {
		ch <- 0.5
	}(ch)
	if len(watermark) == 0 {
		return nil, errors.New("水印字符为空")
	}
	if len(ext) == 0 {
		return nil, errors.New("文件类型为空")
	}
	if strings.EqualFold(ext, ".png") {
		return dealPNG(data, watermark)
	} else if strings.EqualFold(ext, ".jpg") {
		return dealJPG(data, watermark)
	} else if strings.EqualFold(ext, ".jpeg") {
		return dealJPG(data, watermark)
	}
	return nil, errors.New("图片处理失败")
}
func dealPNG(data []byte, watermark string) ([]byte, error) {
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	wmT, err := generateWaterMark(img, watermark, 20)
	if err != nil {
		return nil, err
	}
	marked, err := generateMarkedImage(img, wmT)
	if err != nil {
		return nil, err
	}
	shadowMark, err := generateWaterMark(img, watermark, 5)
	if err != nil {
		return nil, err
	}
	target := insertShadow(marked, shadowMark, "asdasda", 5000)
	buf := new(bytes.Buffer)
	err = png.Encode(buf, target)
	if err != nil {
		return nil, err
	}
	outBytes := buf.Bytes()
	return outBytes, nil
}

func dealJPG(data []byte, watermark string) ([]byte, error) {
	img, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	wmT, err := generateWaterMark(img, watermark, 20)
	if err != nil {
		return nil, err
	}
	marked, err := generateMarkedImage(img, wmT)
	if err != nil {
		return nil, err
	}
	target := insertShadow(marked, wmT, "asdasda", 1000)
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, target, &jpeg.Options{Quality: 100})
	if err != nil {
		return nil, err
	}
	outBytes := buf.Bytes()
	return outBytes, nil
}
func generateWaterMark(ori image.Image, watermark string, rate int) (image.Image, error) {
	font, err := initFont()
	bg := image.Transparent
	fg := image.NewUniform(&color.NRGBA{150, 150, 150, 255})
	bounds := ori.Bounds()
	unitH := bounds.Dy() / rate
	if unitH < minHeight {
		unitH = minHeight
	}
	fmt.Printf("unitH:%d\r\n", unitH)
	width := 0
	for _, c := range watermark {
		fmt.Printf("%c\r\n", c)
		if !isHan(c) {
			fmt.Println("字母")
			width = width + (unitH / 2) + 2
		} else {
			fmt.Println("汉字")
			width = width + unitH
		}
		fmt.Printf("宽度,%d", width)
	}

	wmT := image.NewRGBA(image.Rect(0, 0, width, unitH+10))
	draw.Draw(wmT, wmT.Bounds(), bg, image.ZP, draw.Src)
	context := freetype.NewContext()
	context.SetClip(wmT.Bounds())
	context.SetDst(wmT)
	context.SetSrc(fg)
	context.SetFont(font)
	context.SetFontSize(float64(unitH))
	pt := freetype.Pt(0, (unitH/2)+8+int(context.PointToFixed(float64(unitH))>>8))
	_, err = context.DrawString(watermark, pt)

	return wmT, err
}
func generateMarkedImage(ori image.Image, mark image.Image) (image.Image, error) {
	bounds := ori.Bounds()
	target := image.NewRGBA(bounds)
	offset := image.Pt(ori.Bounds().Dx()-mark.Bounds().Dx()-10, ori.Bounds().Dy()-mark.Bounds().Dy()-10)
	draw.Draw(target, bounds, ori, image.ZP, draw.Src)
	draw.Draw(target, mark.Bounds().Add(offset), mark, image.ZP, draw.Over)
	return target, nil
}
func initFont() (*truetype.Font, error) {
	fontBytes, err := Asset("SiYuanHeiTiJiuZiXing-Regular-2.ttf")
	if err != nil {
		return nil, err
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}
	return font, nil
}

func isHan(r rune) bool {
	return unicode.Is(unicode.Han, r)
}
