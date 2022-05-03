package gomfcc

import (
	"math"
)

func hz2Mel(hz int) float64 {
	return 1127.0 * math.Log(1.0+float64(hz)/700)
}

func mean[T int | float64](ts []T) float64 {
	var sum T
	for _, t := range ts {
		sum += t
	}
	return float64(sum) / float64(len(ts))
}

func matrixT(m [][]float64) [][]float64 {
	mt := make([][]float64, len(m[0]))
	for i := 0; i < len(m[0]); i++ {
		mt[i] = make([]float64, len(m))
	}

	for i := 0; i < len(m); i++ {
		for j := 0; j < len(m[0]); j++ {
			mt[j][i] = m[i][j]
		}
	}
	return mt
}

func matrixDot(x, y [][]float64) [][]float64 {

	ar := len(x)
	ac := len(x[0])
	br := len(y)
	bc := len(y[0])

	// 縦横のサイズが合わない場合
	if ac != br {
		panic("wrong matrix type")
	}

	c := make([][]float64, ar)
	for i := 0; i < ar; i++ {
		c[i] = make([]float64, bc)
		for j := 0; j < bc; j++ {
			for k := 0; k < ac; k++ {
				c[i][j] += x[i][k] * y[k][j]
			}
		}
	}
	return c
}
