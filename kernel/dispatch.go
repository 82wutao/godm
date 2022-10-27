package kernel

import (
	"dm.net/datamine/kernel/types"
)

var procChains = make(map[string]types.ProcessChain)

func RegisterProcess(name string, chain types.ProcessChain) {
	procChains[name] = chain
}

// TODO add a context ,schedule arg,callback
func Dispatch(input *types.DataPresenting) {
	// TODO find a compute-node that is power enough to compute input
	//TODO 上下文 初始化

	consumer := func(src *types.DataPresenting) {

		chain := procChains[src.Name]
		if chain == nil {
			return
		}
		trace := chain.OnStarted(src)

		arg := src
		for i := 0; i < chain.Subprocesses(); i++ {
			sp := chain.Subprocess(i)
			out, e := sp.Process(arg, trace)
			if e != nil {
				go chain.OnError(sp, e, trace)
				return
			}
			arg = out
		}
		combine := chain.GetCombine()
		callback := chain.GetCompleteCallback()
		chain.Reduce(arg, combine, callback, trace)
	}
	go consumer(input)
}
