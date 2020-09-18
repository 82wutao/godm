package util

import "reflect"

// Iterator 迭代器
type Iterator interface {
	More() bool
	Next() interface{}
}

// SliceIterator 切片迭代器
type SliceIterator struct {
	i        int
	elements reflect.Value
}

// More 是否还有更多元素
func (si *SliceIterator) More() bool {
	return si.i < si.elements.Len()
}

// Next 下一个元素
func (si *SliceIterator) Next() interface{} {
	ele := si.elements.Index(si.i)
	ret := ele.Interface()
	si.i++
	return ret
}

// MapPair map的键值对
type MapPair interface {
	Key() interface{}
	Value() interface{}
}
type kvPair struct {
	key   reflect.Value
	value reflect.Value
}

func (kv *kvPair) Key() interface{} {
	return kv.key.Interface()
}
func (kv *kvPair) Value() interface{} {
	return kv.value.Interface()
}

// MapIterator map迭代器
type MapIterator struct {
	iter *reflect.MapIter
}

// More 是否还有更多元素
func (mi *MapIterator) More() bool {
	return mi.iter.Next()
}

// Next 下一个元素
func (mi *MapIterator) Next() interface{} {
	return &kvPair{
		key:   mi.iter.Key(),
		value: mi.iter.Value(),
	}
}

// NewSliceIterator 从Slice构建一个切片迭代器
func NewSliceIterator(slice interface{}) Iterator {
	iter := SliceIterator{
		i:        0,
		elements: reflect.ValueOf(slice),
	}
	return &iter
}

// NewMapIterator 从Map构建一个切片迭代器
func NewMapIterator(mapping interface{}) Iterator {
	iter := MapIterator{
		iter: reflect.ValueOf(mapping).MapRange(),
	}
	return &iter
}
