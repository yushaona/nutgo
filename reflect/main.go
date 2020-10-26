package main

import (
	"fmt"
	"reflect"
)

func main() {
	ExampleReflect()


	
	fmt.Println("----------------------------------")
	var a []byte = make([]byte, 3)
	fmt.Println(a)
	checkType(a)
	fmt.Println(a)

	fmt.Println("----------------------------------")

	fmt.Println(reflect.ValueOf(a).Type())        //  []uint8
	fmt.Println(reflect.ValueOf(a).Type().Elem()) // uint8
	fmt.Println("----------------------------------")
	var b []byte = make([]byte, 3)
	fmt.Println(b)
	checkType(&b)
	fmt.Println(b)
}
