package gomfcc

import "math"

func (m *GoMFCC) makeMelFilterBank(numMelFilterBank, fftSize int) [][]float64 {
	// メル軸での最大周波数
	melHighFreq := hz2Mel(m.option.HighFrequency)
	// メル軸での最小周波数
	melLowFreq := hz2Mel(m.option.LowFrequency)

	// 最小から最大まで等間隔な周波数を得る
	melPoints := make([]float64, 0)
	for p := melLowFreq; p < melHighFreq; p += float64(+2) {
		melPoints = append(melPoints, p)
	}
	melPoints = append(melPoints, melHighFreq)

	// パワースペクトルの次元数
	dimSpectrum := fftSize/2 + 1

	// メルフィルタバンク
	melFilterBank := make([][]float64, numMelFilterBank)
	for i := range melFilterBank {
		melFilterBank[i] = make([]float64, dimSpectrum)
	}
	for i := 0; i < numMelFilterBank; i++ {
		leftMel := melPoints[i]
		centerMel := melPoints[i+1]
		rightMel := melPoints[i+2]
		// パワースペクトルの各ビンに対する重みを計算
		for j := 0; j < dimSpectrum; j++ {
			// 各ビンに対応するHz軸周波数を計算
			freq := 1.0 * j * m.sampleRate / (2.0 * dimSpectrum)
			// メル周波数に変換
			mel := hz2Mel(freq)
			// ビンが三角フィルタの範囲内なら重みを計算
			if leftMel < mel && mel < rightMel {
				if mel <= centerMel {
					melFilterBank[i][j] = (mel - leftMel) / (centerMel - leftMel)
				} else {
					melFilterBank[i][j] = (rightMel - mel) / (rightMel - centerMel)
				}
			}
		}
	}
	return melFilterBank
}

func (m *GoMFCC) makeDCTMatrix(numMelFilterBank, dimMFCC int) [][]float64 {

	dctMatrix := make([][]float64, dimMFCC)
	for i := range dctMatrix {
		dctMatrix[i] = make([]float64, numMelFilterBank)
	}

	// zero
	for i := 0; i < numMelFilterBank; i++ {
		dctMatrix[0][i] = 1.0 / math.Sqrt(float64(numMelFilterBank))
	}

	for i := 1; i < dimMFCC; i++ {
		tmp := float64(i) * math.Pi / (2.0 * float64(numMelFilterBank))
		for j := 0; j < numMelFilterBank; j++ {
			dctMatrix[i][j] = math.Sqrt(2.0/float64(numMelFilterBank)) * math.Cos(float64(2*j+1)*tmp)
		}
	}

	return dctMatrix
}

func (m *GoMFCC) makeLifter(lifterCoef, dimMFCC int) []float64 {
	Q := float64(lifterCoef)
	lifter := make([]float64, dimMFCC)
	for i := 0; i < len(lifter); i++ {
		lifter[i] = 1.0 + 0.5*Q*math.Sin(math.Pi*float64(i)/Q)
	}
	return lifter
}
