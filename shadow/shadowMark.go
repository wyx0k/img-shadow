package shadow

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"io"
	"math"
	"math/rand"

	dsp "github.com/mjibson/go-dsp/fft"
)

func readShadow(ori image.Image) image.Image {
	r, g, b := extractChanel(ori)
	shiftedR := shift(dsp.FFT2(r))
	shiftedG := shift(dsp.FFT2(g))
	shiftedB := shift(dsp.FFT2(b))
	result := applyMatrixByMatrix(shiftedR, plus, shiftedG)
	result = applyMatrixByMatrix(result, plus, shiftedB)
	return convertChanel(ori, shiftedR, shiftedG, shiftedB)
}
func insertShadow(ori image.Image, mark image.Image, secret string, alpha int) image.Image {
	markLayer := fillMarkLayer(extractGray(mark), ori.Bounds())
	alphaMarkLayer := applyMatrix(float64(alpha), multiply, markLayer)
	r, g, b := extractChanel(ori)
	shiftedR := shift(dsp.FFT2(r))
	markedR := applyMatrixByMatrix(shiftedR, plus, alphaMarkLayer)
	shiftedG := shift(dsp.FFT2(g))
	markedG := applyMatrixByMatrix(shiftedG, plus, alphaMarkLayer)
	shiftedB := shift(dsp.FFT2(b))
	markedB := applyMatrixByMatrix(shiftedB, plus, alphaMarkLayer)
	resultR := dsp.IFFT2(ishift(markedR))
	resultG := dsp.IFFT2(ishift(markedG))
	resultB := dsp.IFFT2(ishift(markedB))
	return convertChanel(ori, resultR, resultG, resultB)
}
func fillMarkLayer(matrix [][]complex128, bounds image.Rectangle) [][]complex128 {
	limith := len(matrix)
	limitw := len(matrix[0])

	w := bounds.Dx()
	h := bounds.Dy()
	fmt.Printf("limitw:%d - limith：%d", limitw, limith)
	fmt.Printf("    w:%d - h：%d  ", w, h)
	layer := make([][]complex128, h)
	for i := 0; i < h; i++ {
		if len(layer[i]) == 0 {
			tmp := make([]complex128, w)
			layer[i] = tmp
		}
		if len(layer[h-1-i]) == 0 {
			tmp2 := make([]complex128, w)
			layer[h-1-i] = tmp2
		}
		for j := 0; j < w; j++ {
			if i < limith && j < limitw {
				layer[i][j] = matrix[i][j]
				// fmt.Printf("    i:%d - j：%d  \n", i, j)
				layer[h-1-i][w-1-j] = matrix[i][j]
				// fmt.Printf("	h-1-i:%d - w-1-j：%d  \n", h-1-i, w-1-j)
			}

		}
	}
	return layer
}
func convertChanel(ori image.Image, matrixR, matrixG, matrixB [][]complex128) image.Image {
	bounds := ori.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	result := image.NewNRGBA(bounds)
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			rgba := ori.At(j, i) //第j列第i行，坐标：(j,i)
			_, _, _, a := rgba.RGBA()
			r := real(matrixR[i][j]) / 255
			if r > 255 {
				r = 255
			}
			g := real(matrixG[i][j]) / 255
			if g > 255 {
				g = 255
			}
			b := real(matrixB[i][j]) / 255
			if b > 255 {
				b = 255
			}
			result.Set(j, i, color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}
	return result
}
func convertGray(ori image.Image, matrix [][]complex128) image.Image {
	bounds := ori.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	result := image.NewNRGBA(bounds)
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			g := real(matrix[i][j]) / 255
			if g > 255 {
				g = 0
			} else {
				g = 255
			}
			result.Set(j, i, color.Gray{uint8(g)})
		}
	}
	return result
}
func extractGray(ori image.Image) [][]complex128 {
	bounds := ori.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	result := make([][]complex128, h)
	for i := 0; i < h; i++ {
		line := make([]complex128, w)
		for j := 0; j < w; j++ {
			rgba := ori.At(j, i)
			y := color.GrayModel.Convert(rgba).(color.Gray).Y
			if y > 0 {
				y = 0
			} else {
				y = 255
			}
			line[j] = complex(float64(y), 0)
		}
		result[i] = line
	}
	return result
}
func extractChanel(ori image.Image) (r, g, b [][]complex128) {
	bounds := ori.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	resultR := make([][]complex128, h)
	resultG := make([][]complex128, h)
	resultB := make([][]complex128, h)
	for i := 0; i < h; i++ {
		lineR := make([]complex128, w)
		lineG := make([]complex128, w)
		lineB := make([]complex128, w)
		for j := 0; j < w; j++ {
			rgba := ori.At(j, i)
			r, g, b, _ := rgba.RGBA()
			lineR[j] = complex(float64(r), 0)
			lineG[j] = complex(float64(g), 0)
			lineB[j] = complex(float64(b), 0)
		}
		resultR[i] = lineR
		resultG[i] = lineG
		resultB[i] = lineB
	}
	return resultR, resultG, resultB
}
func applyMatrixByMatrix(matrix1 [][]complex128, f func(a, b float64) float64, matrix2 [][]complex128) [][]complex128 {
	h := len(matrix2)
	w := len(matrix2[0])
	limith := len(matrix1)
	limitw := len(matrix1[0])
	fmt.Println()
	fmt.Printf("limitw:%d - limith：%d \n", limitw, limith)
	fmt.Printf("    w:%d - h：%d  \n", w, h)
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			if h <= limith && w <= limitw {
				matrix1[i][j] = complex(f(real(matrix2[i][j]), real(matrix1[i][j])), imag(matrix1[i][j]))
			}
		}
	}
	return matrix1
}
func applyMatrix(alpha float64, f func(a, b float64) float64, matrix [][]complex128) [][]complex128 {
	h := len(matrix)
	w := len(matrix[0])
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			matrix[i][j] = complex(f(real(matrix[i][j]), alpha), imag(matrix[i][j]))
		}
	}
	return matrix
}
func multiply(a, b float64) float64 {
	return a * b
}
func plus(a, b float64) float64 {
	return a + b
}
func minus(a, b float64) float64 {
	return a - b
}

// func getShadow(bounds image.Rectangle, mark, secret string) image.Image {
// 	width := bounds.Dx()
// 	height := bounds.Dy()
// 	length := math.Min(float64(width), float64(height))
// 	var radius float64
// 	if length <= 100 {
// 		radius = length - 20
// 	} else if length <= 600 {
// 		radius = (length / 2) - 20
// 	} else {
// 		radius = (length / 4) - 20
// 	}
// 	watermark := image.NewNRGBA(bounds)
// 	draw.Draw(watermark, watermark.Bounds(), image.NewUniform(color.Black), image.ZP, draw.Src)
// 	center := image.Point{width / 2, height / 2}
// 	vector := randVector("ellie", 50)
// 	pts := indices(center, int(radius), 50)
// 	for i, pt := range pts {
// 		watermark.SetNRGBA(pt.X, pt.Y, color.NRGBA{0, 0, 0, uint8(255 * vector[i])})
// 	}
// 	return watermark
// }
func randVector(secret string, length int) []int {
	md5gen := md5.New()
	io.WriteString(md5gen, secret)
	seed := binary.BigEndian.Uint64(md5gen.Sum(nil))
	return randVectorByInt(int64(seed), length)
}
func randVectorByInt(secret int64, length int) []int {
	rand.Seed(secret)
	result := make([]int, length)
	for i := 0; i < length; i++ {
		result[i] = rand.Intn(2)
	}
	return result
}

func x(d, cx, radius, vectorLength int) int {
	return cx + int(float64(radius)*math.Cos(float64(d)*2*math.Pi/float64(vectorLength)))
}
func y(d, cy, radius, vectorLength int) int {
	return cy + int(float64(radius)*math.Sin(float64(d)*2*math.Pi/float64(vectorLength)))
}

func indices(
	center image.Point,
	radius, vectorLength int,
) []image.Point {
	cx, cy := center.X, center.Y
	result := make([]image.Point, vectorLength)
	for i := 0; i < vectorLength; i++ {
		result[i] = image.Pt(x(i, cx, radius, vectorLength), y(i, cy, radius, vectorLength))
	}
	return result
}
