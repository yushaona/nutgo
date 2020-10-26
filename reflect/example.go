package main

import (
	"errors"
	"fmt"
	"reflect"
)

// binary.go中摘抄的一段代码  主要是Type和Kind的使用
func intDataSize(data interface{}) int {
	switch data := data.(type) { // 根据数据类型,获取数据大小
	case bool, int8, uint8, *bool, *int8, *uint8:
		return 1
	case []bool:
		return len(data)
	case []int8:
		return len(data)
	case []uint8:
		return len(data)
	case int16, uint16, *int16, *uint16:
		return 2
	case []int16:
		return 2 * len(data)
	case []uint16:
		return 2 * len(data)
	case int32, uint32, *int32, *uint32:
		return 4
	case []int32:
		return 4 * len(data)
	case []uint32:
		return 4 * len(data)
	case int64, uint64, *int64, *uint64:
		return 8
	case []int64:
		return 8 * len(data)
	case []uint64:
		return 8 * len(data)
	case float32, *float32:
		return 4
	case float64, *float64:
		return 8
	case []float32:
		return 4 * len(data)
	case []float64:
		return 8 * len(data)
	}
	return 0
}

//
func checkType(data interface{}) {
	if n := intDataSize(data); n != 0 {
		fmt.Println("dataSize:", n)
		bs := make([]byte, n)
		bs[0] = 0x34
		bs[1] = 0x35
		bs[2] = 0x36
		switch data := data.(type) {
		case *bool:
			*data = bs[0] != 0
		case *int8:
			*data = int8(bs[0])
		case *uint8:
			*data = bs[0]
		case []uint8:
			copy(data, bs)
		default:
			n = 0
		}

		if n != 0 {
			return
		}
	}
	//上面是判断类型,如果没有满足设定了,就走reflect方式
	v := reflect.ValueOf(data)
	fmt.Println(v.Kind())
	fmt.Println(v.Type())
	size := -1
	switch v.Kind() {
	case reflect.Ptr:
		v = v.Elem()
		size = dataSize(v)
	case reflect.Slice:
		size = dataSize(v)
	}

	fmt.Println("size:", size)
	if size < 0 {

		fmt.Println(errors.New("binary.Read: invalid type " + reflect.TypeOf(data).String()))
		return
	}

	return

}

func dataSize(v reflect.Value) int {
	switch v.Kind() {
	case reflect.Slice:
		if s := sizeof(v.Type().Elem()); s >= 0 {
			return s * v.Len()
		}
		return -1

	// case reflect.Struct:
	// 	t := v.Type()
	// 	if size, ok := structSize.Load(t); ok {
	// 		return size.(int)
	// 	}
	// 	size := sizeof(t)
	// 	structSize.Store(t, size)
	// 	return size

	default:
		return sizeof(v.Type())
	}
}

func sizeof(t reflect.Type) int {
	switch t.Kind() {
	case reflect.Array:
		if s := sizeof(t.Elem()); s >= 0 {
			return s * t.Len()
		}

	case reflect.Struct:
		sum := 0
		for i, n := 0, t.NumField(); i < n; i++ {
			s := sizeof(t.Field(i).Type)
			if s < 0 {
				return -1
			}
			sum += s
		}
		return sum

	case reflect.Bool,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return int(t.Size())
	}

	return -1
}
