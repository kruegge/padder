package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"reflect"
	"strconv"
	"time"
	"unsafe"

	"golang.org/x/tools/go/packages"
)

// CustomTypes is a map to store custom type names and their reflect.Type
var CustomTypes = map[string]reflect.Type{
	"MySlice":   reflect.TypeOf(MySlice{}),
	"MyStruct":  reflect.TypeOf(MyStruct{}),
	"time.Time": reflect.TypeOf(time.Time{}),
	"Some":      reflect.TypeOf(Some{}),
}

// ReportStructMemory reports memory usage and padding for a given struct type
func ReportStructMemory(v interface{}) {
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Struct {
		fmt.Println("Error: Not a struct type")
		return
	}

	// Report the unsafe.Sizeof() at the beginning
	fmt.Printf("Unsafe size of struct: %d bytes\n", t.Size())

	fmt.Printf("Analyzing struct: %s\n", t.Name())

	var totalSize uintptr = 0
	var previousOffset uintptr = 0

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		offset := field.Offset // Correct way to get field offset

		// Calculate padding from the previous field
		if i > 0 {
			padding := offset - previousOffset
			if padding > 0 {
				fmt.Printf("  Padding: %d bytes\n", padding)
			}
		}

		fieldSize := field.Type.Size()
		if field.Type.Kind() == reflect.Array {
			fieldSize = uintptr(field.Type.Len()) * field.Type.Elem().Size()
		}

		fmt.Printf("  Field: %-10s Offset: %-3d Size: %-2d Align: %-2d Type: %s\n",
			field.Name, offset, fieldSize, field.Type.Align(), field.Type)

		previousOffset = offset + field.Type.Size()
	}

	// Report total struct size
	totalSize = t.Size()
	fmt.Printf("Total struct size: %d bytes\n", totalSize)
}

// LoadStructFromFile loads a struct definition from a Go file
func LoadStructFromFile(filename, structName string) (interface{}, error) {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filename, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	cfg := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax,
		Fset: fset,
	}
	pkgs, err := packages.Load(cfg, "file="+filename)
	if err != nil {
		return nil, err
	}
	if len(pkgs) == 0 {
		return nil, fmt.Errorf("no packages found in file %s", filename)
	}
	info := pkgs[0].TypesInfo

	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || typeSpec.Name.Name != structName {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			var fields []reflect.StructField
			for _, field := range structType.Fields.List {
				for _, name := range field.Names {
					var tag string
					if field.Tag != nil {
						tag = field.Tag.Value
					}
					fieldType := getType(field.Type, info)
					fields = append(fields, reflect.StructField{
						Name: name.Name,
						Type: fieldType,
						Tag:  reflect.StructTag(tag),
					})
				}
			}

			structValue := reflect.New(reflect.StructOf(fields)).Elem().Interface()
			return structValue, nil
		}
	}

	return nil, fmt.Errorf("struct %s not found in file %s", structName, filename)
}

// getType returns the reflect.Type for a given ast.Expr
func getType(expr ast.Expr, info *types.Info) reflect.Type {
	switch t := expr.(type) {
	case *ast.Ident:
		switch t.Name {
		case "int":
			return reflect.TypeOf(int(0))
		case "int8":
			return reflect.TypeOf(int8(0))
		case "int16":
			return reflect.TypeOf(int16(0))
		case "int32":
			return reflect.TypeOf(int32(0))
		case "int64":
			return reflect.TypeOf(int64(0))

		case "uint":
			return reflect.TypeOf(uint(0))
		case "uint8":
			return reflect.TypeOf(uint8(0))
		case "uint16":
			return reflect.TypeOf(uint16(0))
		case "uint32":
			return reflect.TypeOf(uint32(0))
		case "uint64":
			return reflect.TypeOf(uint64(0))

		case "float32":
			return reflect.TypeOf(float32(0))
		case "float64":
			return reflect.TypeOf(float64(0))

		case "complex64":
			return reflect.TypeOf(complex64(0))
		case "complex128":
			return reflect.TypeOf(complex128(0))

		case "string":
			return reflect.TypeOf("")

		case "bool":
			return reflect.TypeOf(true)

		case "byte":
			return reflect.TypeOf(byte(0)) // Same as uint8
		case "rune":
			return reflect.TypeOf(rune(0)) // Same as int32

		case "uintptr":
			return reflect.TypeOf(uintptr(0))

			// Pointer and unsafe types
		case "unsafe.Pointer":
			return reflect.TypeOf(unsafe.Pointer(nil))
		default:
			// Handle custom types by name
			if typ, ok := CustomTypes[t.Name]; ok {
				return typ
			}
			return reflect.TypeOf(nil)
		}
	case *ast.StarExpr:
		return reflect.PtrTo(getType(t.X, info))
	case *ast.ArrayType:
		elemType := getType(t.Elt, info) // Get the element type
		if t.Len != nil {                // Check if it's a fixed-size array
			lenValue, ok := evalArrayLen(t.Len, info) // Evaluate array length
			if !ok {
				return nil // Handle error or unknown length case
			}
			return reflect.ArrayOf(lenValue, elemType) // Return array type
		}

		return reflect.SliceOf(getType(t.Elt, info))
	case *ast.MapType:
		return reflect.MapOf(getType(t.Key, info), getType(t.Value, info))
	case *ast.SelectorExpr:
		// Handle imported types
		if pkg, ok := t.X.(*ast.Ident); ok {
			typeName := pkg.Name + "." + t.Sel.Name
			if typ, ok := CustomTypes[typeName]; ok {
				return typ
			}
		}
		//if pkg, ok := t.X.(*ast.Ident); ok {
		//	fmt.Println(pkg)
		//	obj := info.Uses[t.Sel]
		//	if obj != nil {
		//		return reflect.TypeOf(obj.Type().(*types.Named).Obj().Name())
		//	}
		//}

		fmt.Printf("Unknown selector type: %s.%s\n", t.X, t.Sel.Name)
		return reflect.TypeOf(nil)
	default:
		fmt.Printf("Unknown expression type: %T\n", expr)
		return reflect.TypeOf(nil)
	}
}

func evalArrayLen(expr ast.Expr, info *types.Info) (int, bool) {
	basicLit, ok := expr.(*ast.BasicLit)
	if !ok || basicLit.Kind != token.INT {
		return 0, false // Not an integer literal
	}
	length, err := strconv.Atoi(basicLit.Value)
	if err != nil {
		return 0, false
	}
	return length, true
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <filename> <structname>")
		return
	}

	filename := os.Args[1]
	structName := os.Args[2]

	structValue, err := LoadStructFromFile(filename, structName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	ReportStructMemory(structValue)
}
