package wave

type avgFilterImpl struct {
	window int
}

// NewAvgFilter build a average filter
func NewAvgFilter(window int) ComplexFiltering {
	return &avgFilterImpl{window: window}
}

// Filtering convert a value in wave to another value
func (avg *avgFilterImpl) Filtering(d float64) float64 {
	return 0
}

// Filtering convert  wave to another wave
func (avg *avgFilterImpl) BatchFiltering(d []float64) []float64 {
	ret := make([]float64, len(d), len(d))
	sum := 0.0

	m := avg.window / 2
	for i := 0; i < m+1; i++ {
		sum += d[0]
	}
	for i := m + 1; i < avg.window; i++ {
		sum += d[i-m-1]
	}
	for i := 0; i < len(d); i++ {
		var rmIdx int
		var addIdx int
		if i-m-1 < 0 {
			rmIdx = 0
		} else {
			rmIdx = i - m - 1
		}
		if i+m >= len(d) {
			addIdx = len(d) - 1
		} else {
			addIdx = i + m
		}

		sum = sum - d[rmIdx]
		sum = sum + d[addIdx]

		ret[i] = sum / float64(avg.window)
	}

	return ret
}
