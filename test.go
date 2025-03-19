package main

import "time"

type MyStruct struct {
	A int
	B string
	C float64
	D bool
	E MySlice
	F *MySlice
	G []MySlice
	H map[string]MySlice
	I []*MySlice
	J time.Time
	K uint32
	L Some
}

type MySlice struct {
	ID   uint64
	Data string
	More bool
}

type MoreStructs struct {
	Data [1000]MyStruct
}

type MoreSlice struct {
	Data [1000]MySlice
}
