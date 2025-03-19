package main

import "time"

type MyStruct struct {
	A int
	B string
	C float64
	D bool
	E SubData
	F *SubData
	G []SubData
	H map[string]SubData
	I []*SubData
	J time.Time
	K uint32
	L Some
}

type SubData struct {
	ID   uint64
	Data string
	More bool
}

type ArrayMyStruct struct {
	Data [1000]MyStruct
}

type ArraySubData struct {
	Data [1000]SubData
}

type ArrayReferenceSubData struct {
	Data *[1000]SubData
}

type SliceSubData struct {
	Data []SubData
}

type SliceReferenceSubData struct {
	Data *[]SubData
}
