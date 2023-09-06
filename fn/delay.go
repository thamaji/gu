package fn

import (
	"reflect"
	"time"
)

func WithDelay[F any](f F, delay time.Duration) F {
	vf := reflect.ValueOf(f)
	typ := vf.Type()
	if typ.Kind() != reflect.Func {
		panic("fn: WithDelay(non-function " + typ.String() + ")")
	}

	var call func([]reflect.Value) []reflect.Value
	if typ.IsVariadic() {
		call = vf.CallSlice
	} else {
		call = vf.Call
	}

	return reflect.MakeFunc(typ, func(args []reflect.Value) []reflect.Value {
		time.Sleep(delay)
		return call(args)
	}).Interface().(F)
}
