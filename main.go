package main

import (
	"fmt"
	"reflect"
	"time"

	coll_list "container/list"

	"dm.net/datamine/common/log"
	"dm.net/datamine/kernel/schedule"
	"dm.net/datamine/kernel/types"
	"dm.net/datamine/visualization"
	// "dm.net/datamine/kernel"
	// "dm.net/datamine/kernel/schedule"
	// "dm.net/datamine/kernel/types"
)

func main() {
	// chain := NewTestChain()
	// kernel.RegisterProcess("test", chain)

	// var data1 = TestRawPresenting{1}
	// var data2 = TestRawPresenting{2}

	// kernel.Dispatch(&types.DataPresenting{Name: "test", Data: &data1})
	// kernel.Dispatch(&types.DataPresenting{Name: "test", Data: &data2})

	// filter := wave.NewMedianFilter(5)

	// ret := filter.BatchFiltering([]float64{9, 5, 2, 10, 7, 4, 1, 6, 3, 8})
	// for i := 0; i < len(ret); i++ {
	// 	fmt.Printf("%f ", ret[i])
	// }
	// fmt.Println()

	// filter = wave.NewAvgFilter(5)

	// ret = filter.BatchFiltering([]float64{689, 5882, 1377, 5828,
	// 	2481, 3037, 5981, 3277, 7732, 9061, 1185, 9374, 2425, 2455, 2453, 4962,
	// 	6144, 6657, 9672, 684, 552, 2224, 7541, 4671, 8147, 9773, 5187, 6035})
	// for i := 0; i < len(ret); i++ {
	// 	fmt.Printf("%f ", ret[i])
	// }
	// fmt.Println()
	logger := log.NewLogger(false, log.DebugLevel,
		[]log.LayoutElement{log.LEVEL, log.DATATIME, log.FILE, log.FUNC, log.LINE, log.MESSAGE},
		log.NewStdoutAppender())
	logger.Error("start visuallization server")

	visualization.AppLaunch()
	select {}
}

func simulatORM(nonnil interface{}, nilp *[]interface{}) {
	ptrValue := reflect.ValueOf(nilp)
	// ptrValue.SetPointer()
	// reflect.AppendSlice()
	sliceValue := ptrValue.Elem()

	n := reflect.Append(sliceValue, reflect.ValueOf(1))
	sliceValue.Set(n)
}

// qps
// tps 进入完成计数一次，一秒内完成成功；在当前秒上累加
type TestChain struct {
	reduceQueue *coll_list.List
	processes   []types.Subprocess
}
type TestRawPresenting struct {
	num int
}
type TestEndPresenting struct {
	sum int
}

func (chain *TestChain) OnStarted(raw *types.DataPresenting) (t *types.Trace) {
	start := time.Now()
	procID := fmt.Sprintf("%s_%d", raw.Name, start.UnixNano())

	return &types.Trace{ID: procID, ReduceNode: nil}
}

func (chain *TestChain) Subprocesses() int {
	return len(chain.processes)
}
func (chain *TestChain) Subprocess(i int) types.Subprocess {
	return chain.processes[i]
}

func (chain *TestChain) OnError(sp types.Subprocess, err error, t *types.Trace) {
	// chain.reduceQueue.Remove(t.ReduceNode)
}

func (chain *TestChain) Reduce(midResultPtr *types.DataPresenting,
	combine interface{}, callback types.OnCompleted, t *types.Trace) {

	runnable := func() {
		raw := midResultPtr.Data.(*TestRawPresenting)
		presenting := combine.(*TestEndPresenting)
		presenting.sum += raw.num

		go callback(presenting)
	}

	schedule.RunConcurrently(runnable)

}
func (tps *TestChain) GetCombine() (combinePtr interface{}) {
	return &TestEndPresenting{0}
}

func (tps *TestChain) GetCompleteCallback() (callback types.OnCompleted) {
	callback = func(combine interface{}) {
		presenting := combine.(*TestEndPresenting)
		fmt.Printf("sum is [%d] \n", presenting.sum)
	}
	return callback
}

func NewTestChain() types.ProcessChain {
	chain := TestChain{
		reduceQueue: coll_list.New(),
		processes:   make([]types.Subprocess, 0, 10),
	}
	return &chain
}

// delay max min mean

// request alive

// err on subprocess and reduce
// slow subprocess, slow reduce
