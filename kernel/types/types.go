package types

import (
	coll_list "container/list"
)

// DataPresenting present data struct
type DataPresenting struct {
	Name string
	Data interface{}
}

type Trace struct {
	ID         string
	ReduceNode *coll_list.Element
}

// OnCompleted is a callback invoked after processing and reduce
type OnCompleted func(resultPtr interface{})

// Subprocess is a interface present a function calling
type Subprocess interface {
	Process(input *DataPresenting, t *Trace) (output *DataPresenting, e error)
	Name() (name string)
}

//TODO 线性串行
// TODO 派生
// TODO 合并
type ProcessChain interface {
	OnStarted(raw *DataPresenting) (t *Trace)
	Subprocesses() int
	Subprocess(i int) Subprocess
	OnError(sp Subprocess, err error, t *Trace)

	//TODO reduce in a concurrently,sync mutex,order queue
	Reduce(midResultPtr *DataPresenting,
		combine interface{}, callback OnCompleted, t *Trace)
	GetCombine() (combinePtr interface{})
	GetCompleteCallback() (callback OnCompleted)
}
