package builtinmodules

import (
	"fmt"
	"time"

	"github.com/LingByte/memole/pkg/object"
)

func newStdModule() *object.ModuleObject {
	exports := map[string]object.Object{
		"nowUnix": &object.NativeFunction{
			Name: "std.nowUnix",
			Fn: func(args []object.Object) object.Object {
				if len(args) != 0 {
					return &object.Null{}
				}
				return &object.Integer{Value: time.Now().Unix()}
			},
		},
		"typeOf": &object.NativeFunction{
			Name: "std.typeOf",
			Fn: func(args []object.Object) object.Object {
				if len(args) != 1 {
					return &object.Null{}
				}
				return &object.String{Value: string(args[0].Type())}
			},
		},
		"print": &object.NativeFunction{
			Name: "std.print",
			Fn: func(args []object.Object) object.Object {
				for i, a := range args {
					if i > 0 {
						fmt.Print(" ")
					}
					fmt.Print(a.Inspect())
				}
				return &object.Null{}
			},
		},
		"println": &object.NativeFunction{
			Name: "std.println",
			Fn: func(args []object.Object) object.Object {
				for i, a := range args {
					if i > 0 {
						fmt.Print(" ")
					}
					fmt.Print(a.Inspect())
				}
				fmt.Println()
				return &object.Null{}
			},
		},
	}

	return &object.ModuleObject{
		Name:    "std",
		Exports: exports,
	}
}
