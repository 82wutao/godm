package metrics

import (
	coll_list "container/list"
	"fmt"
	"time"

	"dm.net/datamine/kernel/schedule"
	"dm.net/datamine/kernel/types"
)

// qps
// tps 进入完成计数一次，一秒内完成成功；在当前秒上累加
type TPSChain struct {
	reduceQueue *coll_list.List
	processes   []types.Subprocess
}
type tpsEndpoint struct {
	start    time.Time
	end      time.Time
	finished bool
}
type TPSPresenting struct {
	timestamp int64
	duration  float64
}

func (tps *TPSChain) OnStarted(raw *types.DataPresenting) (t *types.Trace) {
	start := time.Now()
	procID := fmt.Sprintf("%s_%d", raw.Name, start.UnixNano())
	ele := tps.reduceQueue.PushBack(&tpsEndpoint{start: start, finished: false})

	return &types.Trace{procID, ele}
}
func (tps *TPSChain) Subprocesses() int {
	return len(tps.processes)
}
func (tps *TPSChain) Subprocess(i int) types.Subprocess {
	return tps.processes[i]
}

//TODO concurrent modify list ???
func (tps *TPSChain) OnError(sp types.Subprocess, err error, t *types.Trace) {
	// TODO log error
	tps.reduceQueue.Remove(t.ReduceNode)
}

func (tps *TPSChain) Reduce(midResultPtr *types.DataPresenting,
	combine interface{}, callback types.OnCompleted, t *types.Trace) {

	tpsEndpoin := t.ReduceNode.Value.(*tpsEndpoint)
	tpsEndpoin.finished = true
	tpsEndpoin.end = time.Now()

	runnable := func() {
		for f := tps.reduceQueue.Front(); f != nil; f = tps.reduceQueue.Front() {

			endPoint := f.Value.(*tpsEndpoint)
			if !endPoint.finished {
				break
			}

			tps.reduceQueue.Remove(f)

			dura := endPoint.end.Sub(endPoint.start)
			ts := endPoint.end.Unix()

			go callback(&TPSPresenting{ts, dura.Seconds()})
		}
	}
	schedule.RunConcurrently(runnable)

}
func (tps *TPSChain) GetCombine() (combinePtr interface{}) {
	return nil
}

func (tps *TPSChain) GetCompleteCallback() (callback types.OnCompleted) {
	callback = func(combine interface{}) {
		presenting := combine.(*TPSPresenting)
		fmt.Printf("[%d] expend %f seconds\n", presenting.timestamp, presenting.duration)
	}
	return callback
}

func NewTPSMetricsChain() types.ProcessChain {
	chain := TPSChain{
		reduceQueue: coll_list.New(),
		processes:   make([]types.Subprocess, 0, 10),
	}
	return &chain
}

// delay max min mean

// request alive

// err on subprocess and reduce
// slow subprocess, slow reduce
