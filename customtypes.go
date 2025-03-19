package main

import (
	"reflect"
	"time"
)

// CustomTypes is a map to store custom type names and their reflect.Type - add more types from your files to detect correct sizes
var CustomTypes = map[string]reflect.Type{
	"MySlice":   reflect.TypeOf(MySlice{}),
	"MyStruct":  reflect.TypeOf(MyStruct{}),
	"time.Time": reflect.TypeOf(time.Time{}),
	"Some":      reflect.TypeOf(Some{}),
}
