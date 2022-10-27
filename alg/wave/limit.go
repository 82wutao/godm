package wave

// 限幅滤波法主要用于处理变化较为缓慢的数据，如温度、物体的位置等
type limitingFilterImpl struct {
	max  float64
	last float64
}

func (lf *limitingFilterImpl) Filtering(data float64) float64 {

	if (data-lf.last) > lf.max || (lf.last-data > lf.max) {
		return lf.last
	}

	lf.last = data
	return data

}

// NewLimitingFilter return a interface instance.
// this instance use limiting filtering algorithm
func NewLimitingFilter(maxDiff, init float64) Filtering {
	return &limitingFilterImpl{
		max:  maxDiff,
		last: init,
	}
}
