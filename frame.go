package gomfcc

import (
	"math"
	"math/rand"

	"github.com/yut-kt/gowindow"
)

type frame []float64

func FrameProcessing(
	f []float64,
	ditherCoef float64,
	preEmphasisCoef float64,
) ([]float64, float64) {
	frame := frame(f)
	frame.dithering(ditherCoef)
	frame.dcComponentRemoval()
	logPower := frame.getLogPower()
	frame.preEmphasis(preEmphasisCoef)
	frame.hamming()
	return frame, logPower
}

func (f frame) dithering(coef float64) {
	if coef > 0 {
		// ディザリング
		for n := range f {
			// dither(x(n)) = x(n) + 2Dd(n) - D
			// d(n) = [0.0,1.0)の乱数
			// Dはノイズの大きさを決める自然数
			f[n] += 2.0*coef*rand.Float64() - coef
		}
	}
}

// 直流成分除去
func (f frame) dcComponentRemoval() {
	mean := mean[float64](f)
	for n := range f {
		f[n] -= mean
	}
}

// 対数パワー
func (f frame) getLogPower() float64 {
	sum := 0.0
	for index := range f {
		sum += f[index] * f[index]
	}
	// 対数計算時に-infが出ない様にフロアリング
	if sum < 1e-10 {
		sum = 1e-10
	}
	return math.Log(sum)
}

// pre_emphasis(x(n)) = x(n) - αx(n-1)
func (f frame) preEmphasis(coef float64) {
	for n := len(f) - 1; n > 0; n-- {
		f[n] -= coef * f[n-1]
	}
	f[0] -= coef * f[0]
}

func (f frame) hamming() {
	for index, v := range gowindow.Hamming(f) {
		f[index] = v
	}
}
