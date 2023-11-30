package fn

import (
	"reflect"
	"sync"
)

func WithLock[F any](f F, locker sync.Locker) F {
	vf := reflect.ValueOf(f)
	typ := vf.Type()
	if typ.Kind() != reflect.Func {
		panic("fn: WithLock(non-function " + typ.String() + ")")
	}

	var call func([]reflect.Value) []reflect.Value
	if typ.IsVariadic() {
		call = vf.CallSlice
	} else {
		call = vf.Call
	}

	return reflect.MakeFunc(typ, func(args []reflect.Value) []reflect.Value {
		locker.Lock()
		defer locker.Unlock()
		return call(args)
	}).Interface().(F)
}
