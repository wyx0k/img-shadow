package shadow

import (
	"math"
	"math/cmplx"
)

func fft(signal []complex128) []complex128 {
	n := len(signal)
	if n == 1 {
		return signal
	}
	hl := n / 2
	even := make([]complex128, hl)
	odd := make([]complex128, hl)
	for i := 0; i < hl; i++ {
		even[i] = signal[i*2]
		odd[i] = signal[i*2+1]
	}
	feven := fft(even)
	fodd := fft(odd)
	combined := make([]complex128, n)
	for i := 0; i < hl; i++ {
		combined[i] = feven[i] + omega(n, -i)*fodd[i]
		combined[i+hl] = feven[i] - omega(n, -i)*fodd[i]
	}
	return combined

}

func ifft(data []complex128) []complex128 {
	for i := range data {
		data[i] = cmplx.Conj(data[i])
	}
	scale := complex(float64(len(data)), 0)
	result := fft(data)
	for i := range result {
		result[i] = cmplx.Conj(result[i]) / scale
	}
	return result
}

func fft2d(matrix [][]complex128) [][]complex128 {
	//计算行
	r := len(matrix)
	fftRows := make([][]complex128, r)
	for i := 0; i < r; i++ {
		fftRows[i] = fft(matrix[i])
	}
	//计算列
	transposeRows := transpose(fftRows)
	tr := len(transposeRows)
	tmp := make([][]complex128, tr)
	for i := 0; i < tr; i++ {
		tmp[i] = fft(transposeRows[i])
	}
	result := transpose(tmp)
	return result
}

func ifft2d(matrix [][]complex128) [][]complex128 {
	//计算列
	tmp := transpose(matrix)
	r := len(tmp)
	t1 := make([][]complex128, r)
	for i := 0; i < r; i++ {
		t1[i] = ifft(tmp[i])
	}
	transposeRows := transpose(t1)
	//计算行
	tr := len(transposeRows)
	result := make([][]complex128, tr)
	for i := 0; i < tr; i++ {
		result[i] = ifft(transposeRows[i])
	}
	return result
}

func omega(x, y int) complex128 {
	return cmplx.Exp(complex(0, 2*float64(y)*math.Pi) / complex(float64(x), 0))
}

func transpose(A [][]complex128) [][]complex128 {
	B := make([][]complex128, len(A[0]))
	for i := 0; i < len(A[0]); i++ {
		B[i] = make([]complex128, len(A))
		for j := 0; j < len(A); j++ {
			B[i][j] = A[j][i]
		}
	}
	return B
}

func shift(matrix [][]complex128) [][]complex128 {
	l := len(matrix)
	cIdx := l / 2              //列的交换点
	rIdx := len(matrix[0]) / 2 //行的交换点
	for i := 0; i < l; i++ {
		matrix[i] = shiftArray(matrix[i], rIdx)
	}
	result := transpose(matrix)
	lr := len(result)
	for i := 0; i < lr; i++ {
		result[i] = shiftArray(result[i], cIdx)
	}
	result = transpose(result)
	return result
}
func ishift(matrix [][]complex128) [][]complex128 {
	l := len(matrix)
	cIdx := l / 2              //列的交换点
	rIdx := len(matrix[0]) / 2 //行的交换点
	for i := 0; i < l; i++ {
		matrix[i] = ishiftArray(matrix[i], rIdx)
	}
	result := transpose(matrix)
	lr := len(result)
	for i := 0; i < lr; i++ {
		result[i] = ishiftArray(result[i], cIdx)
	}
	result = transpose(result)
	return result
}

func shiftArray(data []complex128, point int) []complex128 {
	length := len(data)
	tmp := data[length-point:]
	d := data[:length-point]
	return append(tmp, d...)
}
func ishiftArray(data []complex128, point int) []complex128 {
	tmp := data[:point]
	d := data[point:]
	return append(d, tmp...)
}
