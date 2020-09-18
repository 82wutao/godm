package wave

func NewSildingFilter(r int, w []float32) ComplexFiltering {
	if r != len(w) {
		panic("radios must equals w.length")
	}
	return &sildingFilterImpl{
		weights: w,
		radius:  r,
	}
}

type sildingFilterImpl struct {
	weights []float32
	radius  int
}

func (silding *sildingFilterImpl) Filtering(v float64) float64 {
	return 0.0
}

func (silding *sildingFilterImpl) BatchFiltering(data []float64) []float64 {
	ret := make([]float64, len(d), len(d))

	for i := 0; i < len(d); i++ {
		v := float64(0)

		var d float64
		for j := i - (silding.radius - 1); j < i; j++ {
			if j < 0 {
				d = data[0]
			} else {
				d = data[j]
			}
			v += silding.weights[i-j] * d
		}
		for j := i; j < silding.radius-1+i; j++ {
			if j >= len(data) {
				d = data[len(data)-1]
			} else {
				d = data[j]
			}
			v += silding.weights[j-i] * d
		}

		ret[i] = v
	}

	return ret
}
