package fn

import (
	"go/build"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
)

func WithTrace[F any](f F, logger interface{ Println(v ...any) }) F {
	vf := reflect.ValueOf(f)
	typ := vf.Type()
	if typ.Kind() != reflect.Func {
		panic("fn: WithTrace(non-function " + typ.String() + ")")
	}

	var call func([]reflect.Value) []reflect.Value
	if typ.IsVariadic() {
		call = vf.CallSlice
	} else {
		call = vf.Call
	}

	funcName := runtime.FuncForPC(vf.Pointer()).Name()

	return reflect.MakeFunc(typ, func(args []reflect.Value) []reflect.Value {
		_, file, line, _ := runtime.Caller(1)
		for _, gopath := range filepath.SplitList(build.Default.GOPATH) {
			if rel, err := filepath.Rel(gopath, file); err == nil {
				file = rel
				break
			}
		}
		logger.Println("[TRACE]", file+":"+strconv.Itoa(line), funcName)
		return call(args)
	}).Interface().(F)
}
