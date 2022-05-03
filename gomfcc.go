package gomfcc

import (
	"math"

	"github.com/yut-kt/goft"
)

type GoMFCC struct {
	samples    []float64
	sampleRate int
	option     *optionMFCC
}

func NewGoMFCC(samples []float64, sampleRate int, opts ...Option) *GoMFCC {
	option := &optionMFCC{
		LifterCoef:      22,
		DitherCoef:      1.0,
		PreEmphasisCoef: 0.98,
		LowFrequency:    20,
		HighFrequency:   8000,
	}
	for _, opt := range opts {
		opt(option)
	}

	return &GoMFCC{
		samples:    samples,
		sampleRate: sampleRate,
		option:     option,
	}
}

func (m *GoMFCC) GetFeature(
	dimMFCC, numMelFilterBank int,
	numSampleInFrame, numSampleFrameShift int,
) [][]float64 {
	// FFTを行うサンプル数をフレームサイズ以上の2^xに設定
	fftSize := 1
	for fftSize < numSampleInFrame {
		fftSize <<= 2
	}

	melFilterBank := m.makeMelFilterBank(numMelFilterBank, fftSize)
	fbank, logPower := m.getFBANK(numMelFilterBank, fftSize, numSampleInFrame, numSampleFrameShift, melFilterBank)
	dctMatrix := m.makeDCTMatrix(numMelFilterBank, dimMFCC)
	mfcc := matrixDot(fbank, matrixT(dctMatrix))
	lifter := m.makeLifter(m.option.LifterCoef, dimMFCC)
	// リフタリング処理
	for i := range mfcc {
		for j := range mfcc[0] {
			mfcc[i][j] *= lifter[j]
		}
	}
	// MFCCの0次元目を前処理前の波形対数パワーに置き換える
	for i := range mfcc {
		mfcc[i][0] = logPower[i]
	}
	return mfcc
}

// GetFeatureByMS フレーム処理をミリ秒(MS)ごとに行う
func (m *GoMFCC) GetFeatureByMS(
	dimMFCC, numMelFilterBank int,
	msInFrame, msFrameShift float64,
) [][]float64 {
	return m.GetFeature(dimMFCC,
		numMelFilterBank,
		int(float64(m.sampleRate)*msInFrame*0.001),
		int(float64(m.sampleRate)*msFrameShift*0.001),
	)
}

// 対数メルフィルタバンク特徴量と対数パワーを計算
func (m *GoMFCC) getFBANK(
	numFilterBank, fftSize int,
	numSampleInFrame, numSampleFrameShift int,
	melFilterBank [][]float64,
) ([][]float64, []float64) {
	numFrames := (len(m.samples)-numSampleInFrame)/numSampleFrameShift + 1

	fbankFeatures := make([][]float64, numFrames)
	for i := range fbankFeatures {
		fbankFeatures[i] = make([]float64, numFilterBank)
	}
	logPowers := make([]float64, numFrames)

	for frameIndex := 0; frameIndex < numFrames; frameIndex++ {
		startIndex := frameIndex * numSampleFrameShift
		frame := m.samples[startIndex : startIndex+numSampleInFrame]

		var logPow float64
		frame, logPow = FrameProcessing(frame, m.option.DitherCoef, m.option.PreEmphasisCoef)

		// 2^xになるように0埋め
		pad := make([]float64, fftSize)
		for i := 0; i < len(frame); i++ {
			pad[i] = frame[i]
		}
		spectrum, err := goft.FFT(pad)
		if err != nil {
			panic(err)
		}
		spectrum = spectrum[:fftSize/2+1]
		absolute := make([][]float64, 1)
		absolute[0] = make([]float64, len(spectrum))
		for i, s := range spectrum {
			absolute[0][i] = math.Sqrt(math.Pow(real(s), 2) + math.Pow(imag(s), 2))
		}

		fbank := matrixDot(absolute, matrixT(melFilterBank))[0]
		for i, v := range fbank {
			if v < 0.1 {
				fbank[i] = 0.1
			}
			fbank[i] = math.Log(fbank[i])
		}
		fbankFeatures[frameIndex] = fbank
		logPowers[frameIndex] = logPow
	}
	return fbankFeatures, logPowers
}
