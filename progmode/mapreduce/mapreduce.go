package mapreduce

import (
	"reflect"
	"sync"
)

type MappingFunction func(srcSlice interface{}) (dest interface{})
type ReduceFunction func(combine interface{}, destSlice interface{})

type Complete func(interface{})
type Callback func(midResultSlice interface{})
type Defer func(c Callback)

func Map(dataset interface{}, destType reflect.Type, mapFunc MappingFunction, routines int) Defer {
	if dataset == nil {
		// TODO
	}
	dataCollVal := reflect.ValueOf(dataset)
	len := dataCollVal.Len()
	if len == 0 {
		// TODO
	}

	if routines > len {
		routines = len
	}

	d := func(c Callback) {
		// defer function start
		batchSize := len / routines
		batchWave := (len + batchSize - 1) / routines
		resultSet := reflect.MakeSlice(destType, batchWave, batchWave)

		var countdownLatch sync.WaitGroup
		countdownLatch.Add(batchWave)

		for i := 0; i < batchWave; i++ {
			go func(index int) {
				defer countdownLatch.Done()

				sliceVal := dataCollVal.Slice(index*batchSize, (index+1)*batchSize)
				src := sliceVal.Interface()
				dest := mapFunc(src)

				destVal := reflect.ValueOf(dest)
				resultSet.Index(index).Set(destVal)
			}(i)
		}
		countdownLatch.Wait()

		c(resultSet.Interface())
	}
	return d
}

func (d Defer) Reduce(reduce ReduceFunction, combineInput interface{}, callback Complete) {

	c := func(midResultSlice interface{}) {
		midResultVal := reflect.ValueOf(midResultSlice)

		len := midResultVal.Len()
		for i := 0; i < len; i++ {
			resultset := midResultVal.Index(i)
			reduce(combineInput, resultset.Addr().Interface())
		}

		callback(combineInput)
	}
	d(c)
}

/*
 */
