package gomfcc_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/yut-kt/gomfcc"
)

func sampleWave() []float64 {
	const (
		samples = 10000
		tau     = 2 * math.Pi
		rad     = tau / samples
	)

	wave := make([]float64, samples)
	for i := 0; i < samples; i++ {
		wave[i] = math.Sin(rad * float64(i))
	}

	return wave
}

func TestMFCC(t *testing.T) {
	samples := sampleWave()
	sampleRate := 500
	mfcc := gomfcc.NewGoMFCC(samples, sampleRate)
	feature := mfcc.GetFeatureByMS(13, 23, 25, 10)
	fmt.Println("frameLength: ", len(feature))
	fmt.Println("mfcc dims: ", len(feature[0]))
}
