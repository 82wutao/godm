package stream

import (
	"container/list"
	"fmt"
	"reflect"
	"testing"
)

func Test_FilterPipeline(t *testing.T) {
	filter := FilterPipeline{
		result: false,
		function: func(src interface{}) bool {
			ele := src.(int)
			return ele&0x01 == 0
		},
	}

	filter.Accept(0)
	if !filter.NextPipeline() {
		t.Error("filter ,ele:0 match expection,should call next")
	}
	if filter.Result().(int) != 0 {
		t.Error("filter ,ele:0 should be equals 0")
	}
	if filter.TerminateStream() {
		t.Error("filter ,should never terminal")
	}

	filter.Accept(1)
	if filter.NextPipeline() {
		t.Error("filter ,ele:1 does not match expection,should not call next")
	}
	if filter.Result() != nil {
		t.Error("filter ,ele:1 should be equals nil")
	}
	if filter.TerminateStream() {
		t.Error("filter ,should never terminal")
	}

}

func Test_MapPipeline(t *testing.T) {
	mapping := MapPipeline{
		result: nil,
		function: func(src interface{}) interface{} {
			return fmt.Sprintf("%d", src.(int))
		},
	}

	mapping.Accept(0)
	if !mapping.NextPipeline() {
		t.Error("mapping ,should call next forever")
	}
	if mapping.Result().(string) != "0" {
		t.Error("mapping ,should map 0 to '0' ")
	}
	if mapping.TerminateStream() {
		t.Error("filter ,should never terminal")
	}

}

func Test_LimitPipeline(t *testing.T) {
	pipeline := LimitPipeline{
		limit:  3,
		i:      0,
		result: nil,
	}

	collection := make([]int, 0)
	for i := 0; i < 5; i++ {
		pipeline.Accept(i)
		if pipeline.NextPipeline() != (i < 3) {
			t.Errorf("limit ,should call next forever in case that %d < 3", i)
		}
		if pipeline.Result().(int) != i {
			t.Errorf("limit ,%d should equals %d ", pipeline.Result().(int), i)
		}
		if pipeline.TerminateStream() != (i >= 3) {
			t.Error("limit ,should terminal int case that i>=3")
		}

		if pipeline.NextPipeline() {
			r := pipeline.Result().(int)
			collection = append(collection, r)
		}

		if pipeline.TerminateStream() {
			break
		}
	}
	if len(collection) != 3 {
		t.Error("limit ,resultset 's length should be 3")
	}

}

func Test_SkipPipeline(t *testing.T) {
	pipeline := SkipPipeline{
		skip:   3,
		i:      0,
		result: nil,
	}

	collection := make([]int, 0)
	for i := 0; i < 5; i++ {
		pipeline.Accept(i)
		if pipeline.NextPipeline() != (i >= 3) {
			t.Errorf("skip ,should call next forever in case that %d >= 3", i)
		}
		if pipeline.Result().(int) != i {
			t.Errorf("skip ,%d should equals %d ", pipeline.Result().(int), i)
		}
		if pipeline.TerminateStream() {
			t.Error("skip ,should never terminal ")
		}

		if pipeline.NextPipeline() {
			r := pipeline.Result().(int)
			collection = append(collection, r)
		}

		if pipeline.TerminateStream() {
			break
		}
	}
	if len(collection) != 2 {
		t.Error("limit ,resultset 's length should be 3")
	}
}

type person struct {
	age  int
	name string
}

func Test_DistinctPipeline(t *testing.T) {

	hashcode := func(src interface{}) int {
		return src.(*person).age
	}
	equals := func(a, b interface{}) bool {
		pa := a.(*person)
		pb := b.(*person)

		return (pa.age == pb.age) && (pa.name == pb.name)
	}

	pipeline := DistinctPipeline{
		unique: make(map[int][]interface{}),
		hash:   hashcode,
		equals: equals,
		result: nil,
	}

	persons := make([]*person, 5)
	persons[0] = &person{age: 80, name: "a"}
	persons[1] = &person{age: 90, name: "b"}
	persons[2] = &person{age: 80, name: "c"}
	persons[3] = &person{age: 90, name: "b"}
	persons[4] = &person{age: 80, name: "a"}

	collection := make([]interface{}, 0)

	assertFunc := func(p *person, next bool) {
		pipeline.Accept(p)
		if pipeline.NextPipeline() != next {
			t.Errorf("distinct ,calling next expects %t,but actual %t", next, !next)
		}

		var result interface{} = nil
		if next {
			result = p
		}
		if pipeline.Result() != result {
			t.Errorf("distinct ,result expect %v,but %v ", result, pipeline.Result())
		}
		if pipeline.TerminateStream() {
			t.Error("distinct ,should never terminal ")
		}

		if pipeline.NextPipeline() {
			r := pipeline.Result()
			collection = append(collection, r)
		}
	}

	assertFunc(persons[0], true)
	assertFunc(persons[1], true)
	assertFunc(persons[2], true)
	assertFunc(persons[3], false)
	assertFunc(persons[4], false)
	if len(collection) != 3 {
		t.Error("distinct ,resultset 's length should be 3")
	}
}

func Test_ForeachPipeline(t *testing.T) {
	pipeline := ForeachPipeline{
		function: func(src interface{}) {
			fmt.Printf("print a %v\n", src)
		},
	}

	persons := make([]*person, 5)
	persons[0] = &person{age: 80, name: "a"}
	persons[1] = &person{age: 90, name: "b"}
	persons[2] = &person{age: 80, name: "c"}
	persons[3] = &person{age: 90, name: "b"}
	persons[4] = &person{age: 80, name: "a"}

	for _, p := range persons {
		pipeline.Accept(p)
		if !pipeline.NextPipeline() {
			t.Errorf("foreach ,calling next expects true")
		}

		if pipeline.Result() != nil {
			t.Errorf("foreach ,result expect nil ")
		}
		if pipeline.TerminateStream() {
			t.Error("foreach ,should never terminal ")
		}
	}
}

func Test_ReducePipeline(t *testing.T) {
	pipeline := ReducePipeline{
		result: 0,
		accumulator: func(result, src interface{}) interface{} {
			return result.(int) + src.(*person).age
		},
	}

	persons := make([]*person, 5)
	persons[0] = &person{age: 80, name: "a"}
	persons[1] = &person{age: 90, name: "b"}
	persons[2] = &person{age: 80, name: "c"}
	persons[3] = &person{age: 90, name: "b"}
	persons[4] = &person{age: 80, name: "a"}

	for _, p := range persons {
		pipeline.Accept(p)
		if !pipeline.NextPipeline() {
			t.Errorf("reduce ,calling next expects true")
		}

		if pipeline.TerminateStream() {
			t.Error("reduce ,should never terminal ")
		}
	}

	if pipeline.Result().(int) != (80*3 + 90*2) {
		t.Errorf("reduce, sum reduce is wrong")
	}
}

func Test_CollectPipeline(t *testing.T) {
	pipeline := CollectPipeline{
		result: make([]string, 0),
		list:   list.New().Init(),
	}

	persons := make([]*person, 5)
	persons[0] = &person{age: 80, name: "a"}
	persons[1] = &person{age: 90, name: "b"}
	persons[2] = &person{age: 80, name: "c"}
	persons[3] = &person{age: 90, name: "b"}
	persons[4] = &person{age: 80, name: "a"}

	for _, p := range persons {
		pipeline.Accept(p.name)
		if !pipeline.NextPipeline() {
			t.Errorf("collect ,calling next expects true")
		}

		if pipeline.TerminateStream() {
			t.Error("collect ,should never terminal ")
		}
	}

	if len(pipeline.Result().([]string)) != 5 {
		t.Errorf("collect, result should be a slice which of length is 5")
	}
}
func Test_TopKPipeline(t *testing.T) {
	comparator := func(a, b interface{}) int {
		pa := a.(*person)
		pb := b.(*person)

		if pa.age < pb.age {
			return -1
		}
		if pa.age > pb.age {
			return 1
		}
		if pa.name < pb.name {
			return -1
		}
		if pa.name > pb.name {
			return -1
		}

		return 0
	}
	pipeline := TopKPipeline{
		result:     make([]*person, 0),
		k:          3,
		comparator: comparator,
		list:       list.New().Init(),
	}

	persons := make([]*person, 5)
	persons[0] = &person{age: 80, name: "a"}
	persons[1] = &person{age: 90, name: "b"}
	persons[2] = &person{age: 81, name: "c"}
	persons[3] = &person{age: 91, name: "b"}
	persons[4] = &person{age: 80, name: "a"}

	for _, p := range persons {
		pipeline.Accept(p)
		if !pipeline.NextPipeline() {
			t.Errorf("topk ,calling next expects true")
		}

		if pipeline.TerminateStream() {
			t.Error("topk ,should never terminal ")
		}
	}

	topK := pipeline.Result().([]*person)
	if len(topK) != 3 {
		t.Errorf("topk ,lenght expect %d ,but not", pipeline.k)
	}
	if topK[0] != persons[0] {
		t.Errorf("topk ,,top0 does not match expecting, but %v", topK[0])
	}
	if topK[1] != persons[4] {
		t.Errorf("topk ,,top0 does not match expecting, but %v", topK[0])
	}
	if topK[2] != persons[2] {
		t.Errorf("topk ,,top0 does not match expecting, but %v", topK[0])
	}
}

func Test_GroupPipeline(t *testing.T) {

	pipeline := GroupPipeline{
		result:   reflect.ValueOf(make(map[string]*person)),
		groupFuc: func(src interface{}) (interface{}, interface{}) { return src.(*person).name, src },
	}

	persons := make([]*person, 5)
	persons[0] = &person{age: 80, name: "a"}
	persons[1] = &person{age: 90, name: "b"}
	persons[2] = &person{age: 81, name: "c"}
	persons[3] = &person{age: 91, name: "d"}
	persons[4] = &person{age: 80, name: "e"}

	for _, p := range persons {
		pipeline.Accept(p)
		if !pipeline.NextPipeline() {
			t.Errorf("topk ,calling next expects true")
		}

		if pipeline.TerminateStream() {
			t.Error("topk ,should never terminal ")
		}
	}

	mapping := pipeline.Result().(map[string]*person)
	if len(mapping) != 5 {
		t.Errorf("topk ,lenght expect %d ,but not", 5)
	}
}

func Test_GroupAndReducePipeline(t *testing.T) {

	pipeline := GroupAndReducePipeline{
		result:               reflect.ValueOf(make(map[string]int)),
		groupFunc:            func(src interface{}) (interface{}, interface{}) { return src.(*person).name, src },
		reduceResultSupplier: func() interface{} { return int(0) },
		reduceResultAccumulator: func(result, value interface{}) interface{} {
			return result.(int) + value.(*person).age
		},
	}

	persons := make([]*person, 5)
	persons[0] = &person{age: 80, name: "a"}
	persons[1] = &person{age: 90, name: "b"}
	persons[2] = &person{age: 80, name: "a"}
	persons[3] = &person{age: 100, name: "d"}
	persons[4] = &person{age: 110, name: "e"}

	for _, p := range persons {
		pipeline.Accept(p)
		if !pipeline.NextPipeline() {
			t.Errorf("groupreduce ,calling next expects true")
		}

		if pipeline.TerminateStream() {
			t.Error("groupreduce ,should never terminal ")
		}
	}

	mapping := pipeline.Result().(map[string]int)
	if len(mapping) != 4 {
		t.Errorf("groupreduce ,lenght expect %d ,but not", 4)
	}
	if mapping["a"] != 80+80 {
		t.Errorf("groupreduce ,a expect %d ,but not", 160)
	}
	if mapping["b"] != 90 {
		t.Errorf("groupreduce ,b expect %d ,but not", 90)
	}
	if mapping["d"] != 100 {
		t.Errorf("groupreduce ,d expect %d ,but not", 100)
	}
	if mapping["e"] != 110 {
		t.Errorf("groupreduce ,e  expect %d ,but not", 110)
	}
}

// // GroupPipeline 分组
// type GroupAndReducePipeline struct {
// 	result                  reflect.Value
// 	keyGetter               func(src interface{}) interface{}
// 	valueGetter             func(src interface{}) interface{}
// 	reduceResultSupplier    func() interface{}
// 	reduceResultAccumulator func(result, value interface{}) interface{}
// }

// // Result 当前环节的输出结果，也是下一个环节的输入
// func (ppl *GroupAndReducePipeline) Result() interface{} {
// 	return ppl.result.Interface()
// }

// // Accept 当前环节处理输入
// func (ppl *GroupAndReducePipeline) Accept(src interface{}) {

// 	key := ppl.keyGetter(src)
// 	value := ppl.valueGetter(src)

// 	reduceResult := ppl.result.MapIndex(reflect.ValueOf(key))
// 	if reduceResult.IsZero() {
// 		reduceResult = reflect.ValueOf(ppl.reduceResultSupplier())
// 	}

// 	rr := ppl.reduceResultAccumulator(reduceResult.Interface(), value)
// 	ppl.result.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(rr))
// }

// // NextPipeline 是否继续下一个环节
// func (ppl *GroupAndReducePipeline) NextPipeline() bool { return true }

// // TerminateStream 是否短路，是否继续整个流
// func (ppl *GroupAndReducePipeline) TerminateStream() bool { return false }
