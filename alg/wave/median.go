package wave

import "container/list"

//中值滤波比较适用于去掉由偶然因素引起的波动和采样器不稳定而引起的脉动干扰。若被测量值变化比较慢，采用中值滤波法效果会比较好，但如果数据变化比较快，则不宜采用此方法。
//https://www.jianshu.com/p/1937779cbe47
//https://www.cnblogs.com/qsyll0916/p/6691413.html
//https://www.cnblogs.com/yoyo-sincerely/p/6058944.html
//https://blog.csdn.net/nie15870449223/article/details/81252655
type medianFilterImpl struct {
	size int
}

//
func (mf *medianFilterImpl) BatchFiltering(data []float64) []float64 {

	sorted := list.New()
	indexes := list.New()

	medianIndex := func(i int) int {
		if indexes.Len() == mf.size {
			ele := indexes.Front()
			indexes.Remove(ele)

			sortedELe := ele.Value.(*list.Element)
			sorted.Remove(sortedELe)
		}

		n := sorted.Front()
		if n != nil {
			for ; nil != n; n = n.Next() {
				idx := n.Value.(int)
				if data[i] < data[idx] {
					break
				}
			}
		}
		if n == nil {
			ele := sorted.PushBack(i)
			indexes.PushBack(ele)
		} else {
			ele := sorted.InsertBefore(i, n)
			indexes.PushBack(ele)
		}

		median := sorted.Len() / 2
		ele := sorted.Front()
		for i := 0; i < median; i++ {
			ele = ele.Next()
		}
		return ele.Value.(int)
	}

	ret := make([]float64, len(data), len(data))
	m := mf.size / 2
	for i := 0; i < m; i++ {
		medianIndex(0)
	}
	for i := m; i < mf.size-1; i++ {
		medianIndex(i - m)
	}
	for i := 0; i < len(data)-m; i++ {
		mi := medianIndex(i + m)
		ret[i] = data[mi]
	}
	for i := len(data) - m; i < len(data); i++ {
		mi := medianIndex(len(data) - 1)
		ret[i] = data[mi]
	}
	return ret
}

func NewMedianFilter(windowSize int) BatchFiltering {
	return &medianFilterImpl{windowSize}
}
