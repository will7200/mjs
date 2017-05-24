package job

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrParamsNotAdapted = errors.New("The number of params is not adapted.")
)

type Funcs map[string]reflect.Value

func NewFuncs(size int) Funcs {
	return make(Funcs)
}

func (f Funcs) Bind(name string, fn interface{}) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(name + " is not callable.")
		}
	}()
	v := reflect.ValueOf(fn)
	v.Type().NumIn()
	f[name] = v
	return
}

func (f Funcs) Call(name string, params ...interface{}) (result []reflect.Value, err error) {
	if _, ok := f[name]; !ok {
		err = errors.New(name + " does not exist.")
		return
	}
	if len(params) != f[name].Type().NumIn() {
		err = ErrParamsNotAdapted
		return
	}
	in := make([]reflect.Value, len(params))
	//fmt.Println(f[name].Type().In(0))
	for k, param := range params {
		if reflect.TypeOf(param) != f[name].Type().In(k) {
			err = fmt.Errorf("Invalid parameter at index %d, expecting %s (got %s)", k, f[name].Type().In(k), reflect.TypeOf(param))
			return
		}
		in[k] = reflect.ValueOf(param)
	}
	result = f[name].Call(in)
	return
}
