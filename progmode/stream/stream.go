package stream

import (
	"container/list"
	"reflect"

	"dm.net/datamine/progmode/util"
)

type streamImple struct {
	iterator   util.Iterator
	pipelines  []IntermediatePipeline
	terminator TerminalPipeline
}

// NewStreamFromSlice 从切片构建一个流
func NewStreamFromSlice(slice interface{}) Stream {
	return &streamImple{
		iterator:   util.NewIteratorFromSlice(slice),
		pipelines:  make([]IntermediatePipeline, 0),
		terminator: nil,
	}
}

// NewStreamFromMap 从Map构建一个流
func NewStreamFromMap(mapping interface{}) Stream {
	return &streamImple{
		iterator:   util.NewIteratorFromMap(mapping),
		pipelines:  make([]IntermediatePipeline, 0),
		terminator: nil,
	}
}

// Filter 过滤元素
func (s *streamImple) Filter(predicate func(src interface{}) bool) Stream {
	pipeline := &FilterPipeline{
		result:   false,
		function: predicate,
	}
	s.pipelines = append(s.pipelines, pipeline)
	return s
}

// Mapping 映射元素
func (s *streamImple) Mapping(mapFunc func(src interface{}) interface{}) Stream {
	pipeline := &MapPipeline{
		result:   nil,
		function: mapFunc,
	}
	s.pipelines = append(s.pipelines, pipeline)
	return s
}

// Distinct 去重
func (s *streamImple) Distinct(hash func(src interface{}) int, equals func(a, b interface{}) bool) Stream {

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
func (s *streamImple) Limit(l int) Stream {
	pipeline := &LimitPipeline{
		limit:  l,
		i:      0,
		result: nil,
	}
	s.pipelines = append(s.pipelines, pipeline)
	return s
}

// Skip 跳过头部有限的元素
func (s *streamImple) Skip(offset int) Stream {
	pipeline := &SkipPipeline{
		skip:   offset,
		i:      0,
		result: nil,
	}
	s.pipelines = append(s.pipelines, pipeline)
	return s
}

func (s *streamImple) stream() {
	for more := s.iterator.More(); more; more = s.iterator.More() {
		srcEle := s.iterator.Next()

		var lastPipeline IntermediatePipeline = nil
		for _, pipeline := range s.pipelines {
			if lastPipeline != nil {
				pipeline.Accept(lastPipeline.Result())
			} else {
				pipeline.Accept(srcEle)
			}
			lastPipeline = pipeline
			if pipeline.NextPipeline() {
				continue
			}
			break
		}

		currForm := srcEle
		if lastPipeline != nil {
			if !lastPipeline.NextPipeline() {
				continue
			}
			currForm = lastPipeline.Result()
		}

		s.terminator.Accept(currForm)
		if s.terminator.ShortCircuit() {

		}
	}
}

// Foreach 把函数作用在每个元素上
func (s *streamImple) Foreach(f func(src interface{})) {
	terminal := &ForeachPipeline{
		function: f,
	}
	s.terminator = terminal

	s.stream()
}

// Reduce 收敛所有元素为一个形态
func (s *streamImple) Reduce(supplier func() interface{},
	accumulator func(result, ele interface{}) interface{}) interface{} {
	terminal := &ReducePipeline{
		result:      supplier(),
		accumulator: accumulator,
	}
	s.terminator = terminal

	s.stream()
	return terminal.Result()
}

// Collect 收集元素到一个容器
func (s *streamImple) Collect(supplier func() interface{}) interface{} {
	terminal := &CollectPipeline{
		result: supplier(),
		list:   list.New().Init(),
	}
	s.terminator = terminal

	s.stream()
	return terminal.Result()
}

// TopK 有序取头部
func (s *streamImple) TopK(k int, supplier func() interface{}, comparator func(a, b interface{}) int) interface{} {
	terminal := &TopKPipeline{
		result:     supplier(),
		k:          k,
		comparator: comparator,
		list:       list.New().Init(),
	}
	s.terminator = terminal

	s.stream()
	return terminal.Result()
}

// Group 对元素进行分组，一般多用于从slice到map
func (s *streamImple) Group(supplier func() interface{}, groupFuc func(src interface{}) (interface{}, interface{})) interface{} {
	terminal := &GroupPipeline{
		result:   reflect.ValueOf(supplier()),
		groupFuc: groupFuc,
	}
	s.terminator = terminal

	s.stream()
	return terminal.Result()
}

// GroupAndReduce 对元素进行分组，并对同一组的元素进行收敛，
func (s *streamImple) GroupAndReduce(supplier func() interface{}, groupFunc func(src interface{}) (interface{}, interface{}),
	reduceResult func() interface{}, accumulator func(result, value interface{}) interface{}) interface{} {
	terminal := &GroupAndReducePipeline{
		result:                  reflect.ValueOf(supplier()),
		groupFunc:               groupFunc,
		reduceResultSupplier:    reduceResult, //   func() interface{}
		reduceResultAccumulator: accumulator,  // func(result, value interface{}) interface{}
	}
	s.terminator = terminal

	s.stream()
	return terminal.Result()
}
func (s *streamImple) AddIntermediatePipeline(intermediate IntermediatePipeline) Stream {
	s.pipelines = append(s.pipelines, intermediate)
	return s
}
func (s *streamImple) TerminatePipeline(terminator TerminalPipeline) interface{} {
	s.terminator = terminator

	s.stream()
	return terminator.Result()
}

// Stream 流处理
type Stream interface {
	Filter(predicate func(src interface{}) bool) Stream

	// Mapping 映射元素
	Mapping(mapFunc func(src interface{}) interface{}) Stream

	// Distinct 去重
	Distinct(hash func(src interface{}) int, equals func(a, b interface{}) bool) Stream

	// Limit 只取头部有限的元素
	Limit(l int) Stream

	// Skip 跳过头部有限的元素
	Skip(offset int) Stream
	// AddIntermediatePipeline 添加一个中间状态的流水线
	AddIntermediatePipeline(intermediate IntermediatePipeline) Stream

	// Foreach 把函数作用在每个元素上
	Foreach(f func(src interface{}))

	// Reduce 收敛所有元素为一个形态
	Reduce(supplier func() interface{},
		accumulator func(result, ele interface{}) interface{}) interface{}

	// Collect 收集元素到一个容器
	Collect(supplier func() interface{}) interface{}

	// TopK 有序取头部
	TopK(k int, supplier func() interface{}, comparator func(a, b interface{}) int) interface{}

	// Group 对元素进行分组，一般多用于从slice到map
	Group(supplier func() interface{}, groupFuc func(src interface{}) (interface{}, interface{})) interface{}

	// GroupAndReduce 对元素进行分组，并对同一组的元素进行收敛，
	GroupAndReduce(supplier func() interface{}, groupFunc func(src interface{}) (interface{}, interface{}),
		reduceResult func() interface{}, accumulator func(result, value interface{}) interface{}) interface{}

	TerminatePipeline(terminator TerminalPipeline) interface{}
}
