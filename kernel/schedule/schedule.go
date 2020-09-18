package schedule

import (
	"fmt"
	"os"
	"sync"
)

type Runnable func()

func RunConcurrently(r Runnable) {
	go r()
}

var catch = func() {
	if err := recover(); err != nil {
		fmt.Fprint(os.Stdout, err)
	}
}

func RunSync(r Runnable, me *sync.Mutex) {

	go func() {
		me.Lock()
		defer func() {
			me.Unlock()
			catch()
		}()

		r()
	}()
}

func RunInQueue(r Runnable) {
	//TODO add r into a queue consumed by a loop routin
}

// BlockAndWakeup Mutex and cond
type BlockAndWakeup interface {
	Wait4Meeting()
	MeetingAndNotify()
}

func RunWithCondition(r Runnable, cond BlockAndWakeup) {

	go func() {
		cond.Wait4Meeting()
		r()
	}()
}
