package stream

import (
	"container/list"
	"reflect"

	"dm.net/datamine/progmode/util"
)

// Stream 流处理
type Stream struct {
	iterator  util.Iterator
	pipelines []Pipeline
}

// NewSliceStream 从切片构建一个流
func NewSliceStream(slice interface{}) *Stream {
	return &Stream{
		iterator:  util.NewSliceIterator(slice),
		pipelines: make([]Pipeline, 0),
	}
}

// NewMapStream 从Map构建一个流
func NewMapStream(mapping interface{}) *Stream {
	return &Stream{
		iterator:  util.NewMapIterator(mapping),
		pipelines: make([]Pipeline, 0),
	}
}

// Filter 过滤元素
func (s *Stream) Filter(predicate func(src interface{}) bool) *Stream {
	pipeline := &FilterPipeline{
		result:   false,
		function: predicate,
	}
	s.pipelines = append(s.pipelines, pipeline)
	return s
}

// Mapping 映射元素
func (s *Stream) Mapping(mapFunc func(src interface{}) interface{}) *Stream {
	pipeline := &MapPipeline{
		result:   nil,
		function: mapFunc,
	}
	s.pipelines = append(s.pipelines, pipeline)
	return s
}

// Distinct 去重
func (s *Stream) Distinct(hash func(src interface{}) int, equals func(a, b interface{}) bool) *Stream {

	pipeline := &DistinctPipeline{
		unique: make(map[int][]interface{}),
		hash:   hash,
		equals: equals,
		result: nil,
	}
	s.pipelines = append(s.pipelines, pipeline)
	return s
}

// Limit 只取头部有限的元素
func (s *Stream) Limit(l int) *Stream {
	pipeline := &LimitPipeline{
		limit:  l,
		i:      0,
		result: nil,
	}
	s.pipelines = append(s.pipelines, pipeline)
	return s
}

// Skip 跳过头部有限的元素
func (s *Stream) Skip(offset int) *Stream {
	pipeline := &SkipPipeline{
		skip:   offset,
		i:      0,
		result: nil,
	}
	s.pipelines = append(s.pipelines, pipeline)
	return s
}

func (s *Stream) stream() {
	for more := s.iterator.More(); more; more = s.iterator.More() {
		srcEle := s.iterator.Next()

		var lastPipeline Pipeline = nil
		for _, pipeline := range s.pipelines {
			if lastPipeline != nil {
				pipeline.Accept(lastPipeline.Result())
			} else {
				pipeline.Accept(srcEle)
			}
			if pipeline.NextPipeline() {
				lastPipeline = pipeline
				continue
			}

			if pipeline.TerminateStream() {
				return
			}
			break
		}
	}
}

// Foreach 把函数作用在每个元素上
func (s *Stream) Foreach(f func(src interface{})) {
	terminal := &ForeachPipeline{
		function: f,
	}
	s.pipelines = append(s.pipelines, terminal)

	s.stream()
}

// Reduce 收敛所有元素为一个形态
func (s *Stream) Reduce(supplier func() interface{},
	accumulator func(result, ele interface{}) interface{}) interface{} {
	terminal := &ReducePipeline{
		result:      supplier(),
		accumulator: accumulator,
	}
	s.pipelines = append(s.pipelines, terminal)

	s.stream()
	return terminal.Result()
}

// Collect 收集元素到一个容器
func (s *Stream) Collect(supplier func() interface{}) interface{} {
	terminal := &CollectPipeline{
		result: supplier(),
		list:   list.New().Init(),
	}
	s.pipelines = append(s.pipelines, terminal)

	s.stream()
	return terminal.Result()
}

// TopK 有序取头部
func (s *Stream) TopK(k int, supplier func() interface{}, comparator func(a, b interface{}) int) interface{} {
	terminal := &TopKPipeline{
		result:     supplier(),
		k:          k,
		comparator: comparator,
		list:       list.New().Init(),
	}
	s.pipelines = append(s.pipelines, terminal)

	s.stream()
	return terminal.Result()
}

// Group 对元素进行分组，一般多用于从slice到map
func (s *Stream) Group(supplier func() interface{}, groupFuc func(src interface{}) (interface{}, interface{})) interface{} {
	terminal := &GroupPipeline{
		result:   reflect.ValueOf(supplier()),
		groupFuc: groupFuc,
	}
	s.pipelines = append(s.pipelines, terminal)

	s.stream()
	return terminal.Result()
}

// GroupAndReduce 对元素进行分组，并对同一组的元素进行收敛，
func (s *Stream) GroupAndReduce(supplier func() interface{}, groupFunc func(src interface{}) (interface{}, interface{}),
	reduceResult func() interface{}, accumulator func(result, value interface{}) interface{}) interface{} {
	terminal := &GroupAndReducePipeline{
		result:                  reflect.ValueOf(supplier()),
		groupFunc:               groupFunc,
		reduceResultSupplier:    reduceResult, //   func() interface{}
		reduceResultAccumulator: accumulator,  // func(result, value interface{}) interface{}
	}
	s.pipelines = append(s.pipelines, terminal)

	s.stream()
	return terminal.Result()
}

