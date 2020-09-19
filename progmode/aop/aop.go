package aop

import (
	"reflect"
	"runtime"
)

type Joinpointer interface {
	FuncName() string
	FuncProtoType() reflect.Type

	FuncInputs() []reflect.Value
}

// Joinpointer 连接点
type BeforeJoinpointer interface {
	Joinpointer
	SetContext(ctx interface{})
}
type AfterJoinpointer interface {
	Joinpointer
	FuncOutputs() []reflect.Value
	Context() interface{}
}

type AroundJoinpointer interface {
	Joinpointer

	SetFuncInputs([]reflect.Value)
	InvokeFunc() []reflect.Value

	Context() interface{}
	SetContext(ctx interface{})
}
type joinpointerImpl struct {
	name      string
	protoType reflect.Type
	inputs    []reflect.Value
	outputs   []reflect.Value
	function  reflect.Value
	context   interface{}
}

func (jpi *joinpointerImpl) FuncName() string {
	return jpi.name
}
func (jpi *joinpointerImpl) FuncProtoType() reflect.Type {
	return jpi.protoType
}

func (jpi *joinpointerImpl) FuncInputs() []reflect.Value {
	return jpi.inputs
}
func (jpi *joinpointerImpl) SetFuncInputs(args []reflect.Value) {
	jpi.inputs = args
}

func (jpi *joinpointerImpl) InvokeFunc() []reflect.Value {
	return jpi.function.Call(jpi.inputs)
}

func (jpi *joinpointerImpl) FuncOutputs() []reflect.Value {
	return jpi.outputs
}

func (jpi *joinpointerImpl) Context() interface{} {
	return jpi.context
}
func (jpi *joinpointerImpl) SetContext(ctx interface{}) {
	jpi.context = ctx
}

type BeforeAdvice func(jp BeforeJoinpointer)
type AfterAdvice func(jp AfterJoinpointer)
type AroundAdvice func(jp AroundJoinpointer) []reflect.Value

// CreateProxyFunc 创建一个proxy
func CreateProxyFunc(function interface{}, before BeforeAdvice, after AfterAdvice, around AroundAdvice) interface{} {

	funcPrototype := reflect.TypeOf(function)
	proxyDef := func(inputs []reflect.Value) []reflect.Value {

		rawFunc := reflect.ValueOf(function)
		funcName := runtime.FuncForPC(rawFunc.Pointer()).Name()

		jp := joinpointerImpl{
			name:      funcName,
			protoType: funcPrototype,
			inputs:    inputs,
			function:  rawFunc,
		}

		if before != nil {
			before(&jp)
		}

		if around != nil {
			ret := around(&jp)
			jp.outputs = ret
		} else {
			ret := jp.InvokeFunc()
			jp.outputs = ret
		}

		if after != nil {
			after(&jp)
		}

		return jp.outputs
	}

	proxy := reflect.MakeFunc(funcPrototype, proxyDef).Interface()
	return proxy
}
