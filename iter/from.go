package iter

import (
	"bufio"
	"encoding/json"
	"io"
	"io/fs"
	"os"
	"reflect"

	"github.com/thamaji/gu/tuple"
	"golang.org/x/exp/constraints"
)

// 範囲を指定してイテレータをつくる。
func Range[V constraints.Ordered](start V, stop V, step V) Iter[V] {
	cursor := start
	return IterFunc[V](func() (V, bool) {
		if cursor >= stop {
			return *new(V), false
		}
		v := cursor
		cursor += step
		return v, true
	})
}

// 複数の値からイテレータをつくる。
func From[V any](values ...V) Iter[V] {
	cursor := 0
	return IterFunc[V](func() (V, bool) {
		if cursor >= len(values) {
			return *new(V), false
		}
		v := values[cursor]
		cursor++
		return v, true
	})
}

// ポインタからイテレータをつくる。
func FromPtr[V any](p *V) Iter[V] {
	cursor := 0
	return IterFunc[V](func() (V, bool) {
		if p == nil || cursor >= 1 {
			return *new(V), false
		}
		v := *p
		cursor++
		return v, true
	})
}

// スライスからイテレータをつくる。
func FromSlice[V any](slice []V) Iter[V] {
	cursor := 0
	return IterFunc[V](func() (V, bool) {
		if cursor >= len(slice) {
			return *new(V), false
		}
		v := slice[cursor]
		cursor++
		return v, true
	})
}

// マップからイテレータをつくる。
func FromMap[K comparable, V any](m map[K]V) Iter[tuple.T2[K, V]] {
	iter := reflect.ValueOf(m).MapRange()
	next := iter.Next()
	return IterFunc[tuple.T2[K, V]](func() (tuple.T2[K, V], bool) {
		if !next {
			return tuple.NewT2(*new(K), *new(V)), false
		}
		k := iter.Key().Interface().(K)
		v := iter.Value().Interface().(V)
		next = iter.Next()
		return tuple.NewT2(k, v), true
	})
}

// マップのキーからイテレータをつくる。
func FromMapKeys[K comparable, V any](m map[K]V) Iter[K] {
	iter := reflect.ValueOf(m).MapRange()
	next := iter.Next()
	return IterFunc[K](func() (K, bool) {
		if !next {
			return *new(K), false
		}
		k := iter.Key().Interface().(K)
		next = iter.Next()
		return k, true
	})
}

// マップの値からイテレータをつくる。
func FromMapValues[K comparable, V any](m map[K]V) Iter[V] {
	iter := reflect.ValueOf(m).MapRange()
	next := iter.Next()
	return IterFunc[V](func() (V, bool) {
		if !next {
			return *new(V), false
		}
		v := iter.Value().Interface().(V)
		next = iter.Next()
		return v, true
	})
}

// 値と値の有無を受け取ってイテレータをつくる。
func Option[V any](v V, ok bool) Iter[V] {
	if ok {
		return From(v)
	}
	return Empty[V]()
}

// bufio.Scanner からイテレータをつくる。
func FromTextScanner(scanner *bufio.Scanner) Iter[string] {
	next := scanner.Scan()
	return &customIter[string]{
		next: func(ctx Context) (string, bool) {
			if err := scanner.Err(); err != nil {
				ctx.SetErr(err)
			}
			if !next {
				return "", false
			}
			v := scanner.Text()
			next = scanner.Scan()
			return v, true
		},
	}
}

// bufio.Scanner からイテレータをつくる。
func FromBytesScanner(scanner *bufio.Scanner) Iter[[]byte] {
	next := scanner.Scan()
	return &customIter[[]byte]{
		next: func(ctx Context) ([]byte, bool) {
			if err := scanner.Err(); err != nil {
				ctx.SetErr(err)
			}
			if !next {
				return nil, false
			}
			v := scanner.Bytes()
			next = scanner.Scan()
			return v, true
		},
	}
}

// json.Decoder からイテレータをつくる。
func FromJSON[V any](decoder *json.Decoder) Iter[V] {
	more := decoder.More()
	return &customIter[V]{
		next: func(ctx Context) (V, bool) {
			if !more {
				return *new(V), false
			}
			v := *new(V)
			if err := decoder.Decode(&v); err != nil {
				ctx.SetErr(err)
				return v, false
			}
			more = decoder.More()
			return v, true
		},
	}
}

// ディレクトリ内のファイルからイテレータをつくる。
func FromDir(dirname string) Iter[fs.FileInfo] {
	f, err := os.OpenFile(dirname, os.O_RDONLY, 0)
	if err != nil {
		return &customIter[fs.FileInfo]{
			err: err,
			next: func(ctx Context) (fs.FileInfo, bool) {
				return nil, false
			},
		}
	}
	return &customIter[fs.FileInfo]{
		next: func(ctx Context) (fs.FileInfo, bool) {
			list, err := f.Readdir(1)
			if err != nil {
				_ = f.Close()
				if err != io.EOF {
					ctx.SetErr(err)
				}
				return nil, false
			}
			return list[0], false
		},
	}
}

// ディレクトリ内のファイル名からイテレータをつくる。
func FromDirname(dirname string) Iter[string] {
	f, err := os.OpenFile(dirname, os.O_RDONLY, 0)
	if err != nil {
		return &customIter[string]{
			err: err,
			next: func(ctx Context) (string, bool) {
				return "", false
			},
		}
	}
	return &customIter[string]{
		next: func(ctx Context) (string, bool) {
			list, err := f.Readdirnames(1)
			if err != nil {
				_ = f.Close()
				if err != io.EOF {
					ctx.SetErr(err)
				}
				return "", false
			}
			return list[0], false
		},
	}
}

// 関数からイテレータをつくる。
func FromFunc[V any](f func(Context) (V, bool)) Iter[V] {
	return &customIter[V]{
		next: f,
	}
}

type Context interface {
	SetErr(error)
	Err() error
}

type customIter[V any] struct {
	next func(Context) (V, bool)
	err  error
}

func (iter *customIter[V]) Next() (V, bool) {
	if iter.err != nil {
		return *new(V), false
	}
	return iter.next(iter)
}

func (iter *customIter[V]) SetErr(err error) {
	iter.err = err
}

func (iter *customIter[V]) Err() error {
	return iter.err
}
