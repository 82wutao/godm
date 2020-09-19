package aop

import (
	"reflect"
	"runtime"
)

// Joinpointer 连接点
type Joinpointer interface {
	Name() string
	Inputs() []reflect.Value
	Outputs() []reflect.Value
	Invoke() []reflect.Value
	Context() interface{}
	SetContext(ctx interface{})
}
type joinpointerImpl struct {
	name     string
	inputs   []reflect.Value
	outputs  []reflect.Value
	function reflect.Value
	context  interface{}
}

func (jpi *joinpointerImpl) Name() string {
	return jpi.name
}
func (jpi *joinpointerImpl) Inputs() []reflect.Value {
	return jpi.inputs
}
func (jpi *joinpointerImpl) Invoke() []reflect.Value {
	outputs := jpi.function.Call(jpi.inputs)
	jpi.outputs = outputs
	return outputs
}
func (jpi *joinpointerImpl) Outputs() []reflect.Value {
	return jpi.outputs
}
func (jpi *joinpointerImpl) Context() interface{} {
	return jpi.context
}
func (jpi *joinpointerImpl) SetContext(ctx interface{}) {
	jpi.context = ctx
}

type BeforeAdvice func(jp Joinpointer)
type AfterAdvice func(jp Joinpointer)
type AroundAdvice func(jp Joinpointer) []reflect.Value

// CreateProxyFunc 创建一个proxy
func CreateProxyFunc(function interface{}, before BeforeAdvice, after AfterAdvice, around AroundAdvice) interface{} {

	proxyDef := func(inputs []reflect.Value) []reflect.Value {

		rawFunc := reflect.ValueOf(function)
		funcName := runtime.FuncForPC(rawFunc.Pointer()).Name()

		jp := joinpointerImpl{
			name:     funcName,
			inputs:   inputs,
			function: rawFunc,
		}

		if before != nil {
			before(&jp)
		}

		if around != nil {
			ret := around(&jp)
			jp.outputs = ret
		} else {
			ret := jp.Invoke()
			jp.outputs = ret
		}

		if after != nil {
			after(&jp)
		}

		return jp.outputs
	}

	funcPrototype := reflect.TypeOf(function)
	proxy := reflect.MakeFunc(funcPrototype, proxyDef).Interface()
	return proxy
}
