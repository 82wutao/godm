package util

import (
	"testing"
)

func Test_SliceIterator(t *testing.T) {
	sliceInst := make([]string, 0)
	sliceIter := NewSliceIterator(sliceInst)
	if sliceIter == nil {
		t.Error("New slice iterator err")
	}

	more := sliceIter.More()
	if more {
		t.Error("a empty slice should not more elements")
	}

	sliceInst = append(sliceInst, "zero")
	sliceIter = NewSliceIterator(sliceInst)
	if sliceIter == nil {
		t.Error("New slice iterator err")
	}
	more = sliceIter.More()
	if !more {
		t.Error("a slice contain one ele should more elements")
	}
	ele := sliceIter.Next()
	if ele == nil {
		t.Error("a slice containing one ele ,here not found the ele")
	}
	if ele.(string) != "zero" {
		t.Error("a slice containing one ele ,the value of that ele not equals")
	}

	more = sliceIter.More()
	if more {
		t.Error("a slice contain one ele should not more elements")
	}
}

func Test_MapIterator(t *testing.T) {
	mapInst := make(map[int]string)
	mapIter := NewMapIterator(mapInst)
	if mapIter == nil {
		t.Error("New map iterator err")
	}

	more := mapIter.More()
	if more {
		t.Error("a empty map should not more elements")
	}

	mapInst[0] = "zero"
	mapIter = NewMapIterator(mapInst)
	if mapIter == nil {
		t.Error("New map iterator err")
	}
	more = mapIter.More()
	if !more {
		t.Error("a map contain one kv-pair should more elements")
	}
	pair := mapIter.Next()
	if pair == nil {
		t.Error("a map containing one kv-pair ,here not found the kv-pair")
	}
	if pair.(MapPair).Key().(int) != 0 {
		t.Error("a map containing one kv-pair ,the key of that kv-pair not equals")
	}
	if pair.(MapPair).Value().(string) != "zero" {
		t.Error("a map containing one kv-pair ,the value of that kv-pair not equals")
	}

	more = mapIter.More()
	if more {
		t.Error("a map contain one kv-pair should not more elements")
	}
}
