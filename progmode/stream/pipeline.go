package stream

import (
	"container/list"
	"reflect"
)

// Pipeline 流水线环节抽象
type Pipeline interface {
	Result() interface{}
	Accept(src interface{})
}
type IntermediatePipeline interface {
	Pipeline
	NextPipeline() bool
}
type TerminalPipeline interface {
	Pipeline
	ShortCircuit() bool
}

// FilterPipeline 过滤环节
type FilterPipeline struct {
	result   interface{}
	function func(src interface{}) bool
}

// Result 当前环节的输出结果，也是下一个环节的输入
func (fp *FilterPipeline) Result() interface{} { return fp.result }

// Accept 当前环节处理输入
func (fp *FilterPipeline) Accept(src interface{}) {
	fp.result = nil
	if fp.function(src) {
		fp.result = src
	}
}

// NextPipeline 是否继续下一个环节
func (fp *FilterPipeline) NextPipeline() bool { return fp.result != nil }

// MapPipeline 映射环节
type MapPipeline struct {
	result   interface{}
	function func(src interface{}) interface{}
}

// Result 当前环节的输出结果，也是下一个环节的输入
func (mp *MapPipeline) Result() interface{} { return mp.result }

// Accept 当前环节处理输入
func (mp *MapPipeline) Accept(src interface{}) { mp.result = mp.function(src) }

// NextPipeline 是否继续下一个环节
func (mp *MapPipeline) NextPipeline() bool { return true }

// LimitPipeline 头部有限的
type LimitPipeline struct {
	limit  int
	i      int
	result interface{}
}

// Result 当前环节的输出结果，也是下一个环节的输入
func (lp *LimitPipeline) Result() interface{} { return lp.result }

// Accept 当前环节处理输入
func (lp *LimitPipeline) Accept(src interface{}) {
	lp.result = src
	lp.i++
}

// NextPipeline 是否继续下一个环节
func (lp *LimitPipeline) NextPipeline() bool { return lp.i <= lp.limit }

// SkipPipeline 跳过头部
type SkipPipeline struct {
	skip   int
	i      int
	result interface{}
}

// Result 当前环节的输出结果，也是下一个环节的输入
func (sp *SkipPipeline) Result() interface{} { return sp.result }

// Accept 当前环节处理输入
func (sp *SkipPipeline) Accept(src interface{}) {
	sp.result = src
	sp.i++
}

// NextPipeline 是否继续下一个环节
func (sp *SkipPipeline) NextPipeline() bool { return sp.i > sp.skip }

// DistinctPipeline 去重
type DistinctPipeline struct {
	unique map[int][]interface{}
	hash   func(src interface{}) int
	equals func(a, b interface{}) bool
	result interface{}
}

// Result 当前环节的输出结果，也是下一个环节的输入
func (dp *DistinctPipeline) Result() interface{} { return dp.result }

// Accept 当前环节处理输入
func (dp *DistinctPipeline) Accept(src interface{}) {
	dp.result = nil

	hashVal := dp.hash(src)
	collection, existed := dp.unique[hashVal]
	if !existed {
		collection = make([]interface{}, 0)
		dp.unique[hashVal] = collection
	}
	for _, ele := range collection {
		if !dp.equals(src, ele) {
			continue
		}
		return
	}
	dp.result = src
	collection = append(collection, src)
	dp.unique[hashVal] = collection

}

// NextPipeline 是否继续下一个环节
func (dp *DistinctPipeline) NextPipeline() bool { return dp.result != nil }

// 终结系列pipeline

// ForeachPipeline 把函数作用在每个元素上
type ForeachPipeline struct {
	function func(src interface{})
}

// Result 当前环节的输出结果，也是下一个环节的输入
func (foreach *ForeachPipeline) Result() interface{} { return nil }

// Accept 当前环节处理输入
func (foreach *ForeachPipeline) Accept(src interface{}) {
	foreach.function(src)
}

// ShortCircuit 是否短路，是否继续整个流
func (foreach *ForeachPipeline) ShortCircuit() bool { return false }

// ReducePipeline 收敛所有元素为一个形态
type ReducePipeline struct {
	result      interface{}
	accumulator func(result, ele interface{}) interface{}
}

// Result 当前环节的输出结果，也是下一个环节的输入
func (rp *ReducePipeline) Result() interface{} { return rp.result }

// Accept 当前环节处理输入
func (rp *ReducePipeline) Accept(src interface{}) {
	rp.result = rp.accumulator(rp.result, src)
}

// ShortCircuit 是否短路，是否继续整个流
func (rp *ReducePipeline) ShortCircuit() bool { return false }

// CollectPipeline 收集元素
type CollectPipeline struct {
	result interface{}
	list   *list.List
}

// Result 当前环节的输出结果，也是下一个环节的输入
func (cp *CollectPipeline) Result() interface{} {
	temp := reflect.ValueOf(cp.result)

	for ele := cp.list.Front(); ele != nil; ele = ele.Next() {
		temp = reflect.Append(temp, reflect.ValueOf(ele.Value))
	}
	return temp.Interface()
}

// Accept 当前环节处理输入
func (cp *CollectPipeline) Accept(src interface{}) {
	cp.list.PushBack(src)
}

// ShortCircuit 是否短路，是否继续整个流
func (cp *CollectPipeline) ShortCircuit() bool { return false }

// TopKPipeline 有序取头部
type TopKPipeline struct {
	result     interface{}
	k          int
	comparator func(a, b interface{}) int
	list       *list.List
}

// Result 当前环节的输出结果，也是下一个环节的输入
func (ppl *TopKPipeline) Result() interface{} {
	resultContain := reflect.ValueOf(ppl.result)

	i := 0
	for p := ppl.list.Front(); i < ppl.k && p != nil; p = p.Next() {
		resultContain = reflect.Append(resultContain, reflect.ValueOf(p.Value))
		i++
	}

	return resultContain.Interface()
}

// Accept 当前环节处理输入
func (ppl *TopKPipeline) Accept(src interface{}) {

	for p := ppl.list.Front(); p != nil; p = p.Next() {
		f := ppl.comparator(src, p.Value)
		if f < 0 {
			ppl.list.InsertBefore(src, p)
			return
		}
	}
	ppl.list.PushBack(src)
}

// ShortCircuit 是否短路，是否继续整个流
func (ppl *TopKPipeline) ShortCircuit() bool { return false }

// GroupPipeline 分组
type GroupPipeline struct {
	result   reflect.Value
	groupFuc func(src interface{}) (interface{}, interface{})
}

// Result 当前环节的输出结果，也是下一个环节的输入
func (ppl *GroupPipeline) Result() interface{} {
	return ppl.result.Interface()
}

// Accept 当前环节处理输入
func (ppl *GroupPipeline) Accept(src interface{}) {

	key, value := ppl.groupFuc(src)

	ppl.result.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(value))
}

// ShortCircuit 是否短路，是否继续整个流
func (ppl *GroupPipeline) ShortCircuit() bool { return false }

// supplier func() interface{}, keyGetter, valueGetter func(src interface{},
// 	reduceResult func() interface{}, accumulator func(result, value interface{})

// GroupAndReducePipeline 分组
type GroupAndReducePipeline struct {
	result                  reflect.Value
	groupFunc               func(src interface{}) (interface{}, interface{})
	reduceResultSupplier    func() interface{}
	reduceResultAccumulator func(result, value interface{}) interface{}
}

// Result 当前环节的输出结果，也是下一个环节的输入
func (ppl *GroupAndReducePipeline) Result() interface{} {
	return ppl.result.Interface()
}

// Accept 当前环节处理输入
func (ppl *GroupAndReducePipeline) Accept(src interface{}) {

	key, value := ppl.groupFunc(src)

	reduceResult := ppl.result.MapIndex(reflect.ValueOf(key))
	if !reduceResult.IsValid() {
		reduceResult = reflect.ValueOf(ppl.reduceResultSupplier())
	}

	rr := ppl.reduceResultAccumulator(reduceResult.Interface(), value)
	ppl.result.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(rr))
}

// ShortCircuit 是否短路，是否继续整个流
func (ppl *GroupAndReducePipeline) ShortCircuit() bool { return false }
