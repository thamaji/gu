package fn

import (
	"reflect"
	"time"
)

func WithThrottle[F any](f F, delay time.Duration) F {
	vf := reflect.ValueOf(f)
	typ := vf.Type()
	if typ.Kind() != reflect.Func {
		panic("fn: WithThrottle(non-function " + typ.String() + ")")
	}

	var call func([]reflect.Value) []reflect.Value
	if typ.IsVariadic() {
		call = vf.CallSlice
	} else {
		call = vf.Call
	}

	var prev time.Time
	return reflect.MakeFunc(typ, func(args []reflect.Value) []reflect.Value {
		d := delay - time.Since(prev)
		if d > 0 {
			time.Sleep(d)
		}
		prev = time.Now()
		return call(args)
	}).Interface().(F)
}
