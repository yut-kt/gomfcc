package gomfcc

type optionMFCC struct {
	LifterCoef      int
	DitherCoef      float64
	PreEmphasisCoef float64
	LowFrequency    int
	HighFrequency   int
}

type Option func(*optionMFCC)

func LifterCoef(c int) Option {
	return func(mfccOption *optionMFCC) {
		mfccOption.LifterCoef = c
	}
}

func DitherCoef(c float64) Option {
	return func(mfccOption *optionMFCC) {
		mfccOption.DitherCoef = c
	}
}

func PreEmphasisCoef(c float64) Option {
	return func(mfccOption *optionMFCC) {
		mfccOption.PreEmphasisCoef = c
	}
}

func HighFrequency(highFQ int) Option {
	return func(mfccOption *optionMFCC) {
		mfccOption.HighFrequency = highFQ
	}
}

func LowFrequency(lowFQ int) Option {
	return func(mfccOption *optionMFCC) {
		mfccOption.LowFrequency = lowFQ
	}
}
