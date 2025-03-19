# Padder

Padder is a Go project that analyzes the memory usage and padding of struct types. It provides detailed information about the memory layout of structs, including field offsets, sizes, and alignment.

## Installation

To install the project, clone the repository and navigate to the project directory:

```sh
git clone https://github.com/kruegge/padder.git
cd padder
```

Then, build the project using the `go` tool:

```sh
go build
```

This will create an executable named `padder` in the project directory.

## Usage

To analyze a struct type, run the `padder` executable with the name of the struct type as an argument. For example, to analyze the `Person` struct:

```sh
./padder test.go MyStruct
```

This will output detailed information about the memory layout of the `MyStruct` type, including field offsets, sizes, and alignment.

```
Unsafe size of struct: 192 bytes
Analyzing struct: 
  Field: A          Offset: 0   Size: 8  Align: 8  Type: int
  Field: B          Offset: 8   Size: 16 Align: 8  Type: string
  Field: C          Offset: 24  Size: 8  Align: 8  Type: float64
  Field: D          Offset: 32  Size: 1  Align: 1  Type: bool
  Padding: 7 bytes
  Field: E          Offset: 40  Size: 32 Align: 8  Type: main.SubData
  Field: F          Offset: 72  Size: 8  Align: 8  Type: *main.SubData
  Field: G          Offset: 80  Size: 24 Align: 8  Type: []main.SubData
  Field: H          Offset: 104 Size: 8  Align: 8  Type: map[string]main.SubData
  Field: I          Offset: 112 Size: 24 Align: 8  Type: []*main.SubData
  Field: J          Offset: 136 Size: 24 Align: 8  Type: time.Time
  Field: K          Offset: 160 Size: 4  Align: 4  Type: uint32
  Padding: 4 bytes
  Field: L          Offset: 168 Size: 24 Align: 8  Type: main.Some
Total struct size: 192 bytes

```

```sh
./padder test.go ArrayMyStruct
```

```
Unsafe size of struct: 192000 bytes
Analyzing struct: 
  Field: Data       Offset: 0   Size: 192000 Align: 8  Type: [1000]main.MyStruct
Total struct size: 192000 bytes
```
