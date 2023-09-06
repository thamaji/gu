package fn

import (
	"reflect"
)

func Once[F any](f F) F {
	vf := reflect.ValueOf(f)
	typ := vf.Type()
	if typ.Kind() != reflect.Func {
		panic("fn: Once(non-function " + typ.String() + ")")
	}

	var call func([]reflect.Value) []reflect.Value
	if typ.IsVariadic() {
		call = vf.CallSlice
	} else {
		call = vf.Call
	}

	cached := false
	var results []reflect.Value
	return reflect.MakeFunc(typ, func(args []reflect.Value) []reflect.Value {
		if cached {
			return results
		}
		results = call(args)
		cached = true
		return results
	}).Interface().(F)
}
