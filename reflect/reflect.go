package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

//类型和种类的区别范类

type myStruct struct { //类型
}

type newInt int

const (
	m newInt = 3
	n int    = 3
)

/*
	主要的区别在与用type关键字定义的变量名...其他情况下都是一样的

	用type关键字定义了变量,类型和类型名都会变,种类还是一种int

	//比如:狗/猫 都是种类,每种狗,又分为 哈士奇/京巴/柯基等等

*/

func ExampleReflect() {
	diffTypeAndKind() //Type()和Kind()的对比
	structType()      //访问结构体中个的字段类型
	reflect2Value()   //从reflect.ValueOf中获取原始数值
	structValueOf()   // 访问结构体中字段值
	NilAndValild()    // 判断是否为nil或空值
	setValue()        //可寻址的变量才能设置值
	newInstance()     //基于TypeOf类型,新建实例
	call()            //函数调用

	getPtrSize()
}

// type和kind的主要区别
func diffTypeAndKind() {
	fmt.Println("---------diffTypeAndKind----------")
	//自定义结构体
	a := myStruct{}
	fmt.Println(reflect.TypeOf(a))        // main.myStruct
	fmt.Println(reflect.TypeOf(a).Name()) //myStruct
	fmt.Println(reflect.TypeOf(a).Kind()) //struct

	b := &myStruct{}
	fmt.Println(reflect.TypeOf(b))        // *main.myStruct
	fmt.Println(reflect.TypeOf(b).Name()) //空
	fmt.Println(reflect.TypeOf(b).Kind()) //ptr

	typeOfElem := reflect.TypeOf(b).Elem() //解引用
	fmt.Println(typeOfElem)                //  main.myStruct
	fmt.Println(typeOfElem.Name())         //myStruct
	fmt.Println(typeOfElem.Kind())         //struct

	//普通类型
	fmt.Println(reflect.TypeOf(m))        // main.newInt
	fmt.Println(reflect.TypeOf(m).Name()) // newInt
	fmt.Println(reflect.TypeOf(m).Kind()) // int

	fmt.Println(reflect.TypeOf(n))        //int
	fmt.Println(reflect.TypeOf(n).Name()) //int
	fmt.Println(reflect.TypeOf(n).Kind()) //int

}

// 声明一个空结构体
type cat struct {
	Name string
	// 带有结构体tag的字段
	Type int `json:"type" id:"100"`
}

// 如何访问结构体字段的(类型信息)
func structType() {

	/*
		Field(i int) StructField	根据索引，返回索引对应的结构体字段的信息。当值不是结构体或索引超界时发生宕机
		NumField() int	返回结构体成员字段数量。当类型不是结构体或索引超界时发生宕机
		FieldByName(name string) (StructField, bool)	根据给定字符串返回字符串对应的结构体字段的信息。没有找到时 bool 返回 false，当类型不是结构体或索引超界时发生宕机
		FieldByIndex(index []int) StructField	多层成员访问时，根据 []int 提供的每个结构体的字段索引，返回字段的信息。没有找到时返回零值。当类型不是结构体或索引超界时 发生宕机
		FieldByNameFunc( match func(string) bool) (StructField,bool)	根据匹配函数匹配需要的字段。当值不是结构体或索引超界时发生宕机

	*/
	fmt.Println("---------structType----------")
	// 创建cat的实例
	ins := cat{Name: "mimi", Type: 1}
	// 获取结构体实例的反射类型对象
	typeOfCat := reflect.TypeOf(ins)
	// 遍历结构体所有成员
	for i := 0; i < typeOfCat.NumField(); i++ {
		// 获取每个成员的结构体字段类型
		fieldType := typeOfCat.Field(i)
		// 输出成员名和tag
		fmt.Printf("name: %v  tag: '%v'\n", fieldType.Name, fieldType.Tag)
	}
	// 通过字段名, 找到字段类型信息
	if catType, ok := typeOfCat.FieldByName("Type"); ok {
		// 从tag中取出需要的tag
		fmt.Println(catType.Tag.Get("json"), catType.Tag.Get("id"))
	}
}

func reflect2Value() {
	fmt.Println("---------reflect2Value----------")
	// 声明整型变量a并赋初值
	var a int = 1024
	// 获取变量a的反射值对象
	valueOfA := reflect.ValueOf(a)
	fmt.Println("valueOfA.CanSet()", valueOfA.CanSet()) //只有指针类型才能set
	// 获取interface{}类型的值, 通过类型断言转换
	var getA int = valueOfA.Interface().(int)
	// 获取64位的值, 强制类型转换为int类型
	var getA2 int = int(valueOfA.Int())
	fmt.Println(getA, getA2)
}

// 定义结构体
type dummy struct {
	a int
	b string
	// 嵌入字段
	float32
	bool
	next *dummy
}

func structValueOf() {
	fmt.Println("---------structValueOf----------")
	/*
		Field(i int) Value	根据索引，返回索引对应的结构体成员字段的反射值对象。当值不是结构体或索引超界时发生宕机
		NumField() int	返回结构体成员字段数量。当值不是结构体或索引超界时发生宕机
		FieldByName(name string) Value	根据给定字符串返回字符串对应的结构体字段。没有找到时返回零值，当值不是结构体或索引超界时发生宕机
		FieldByIndex(index []int) Value	多层成员访问时，根据 []int 提供的每个结构体的字段索引，返回字段的值。 没有找到时返回零值，当值不是结构体或索引超界时发生宕机
		FieldByNameFunc(match func(string) bool) Value	根据匹配函数匹配需要的字段。找到时返回零值，当值不是结构体或索引超界时发生宕机
	*/

	// 值包装结构体
	d := reflect.ValueOf(dummy{
		next: &dummy{},
	})
	// 获取字段数量
	fmt.Println("NumField", d.NumField())
	// 获取索引为2的字段(float32字段)  索引从0开始
	floatField := d.Field(2)
	// 输出字段类型
	fmt.Println("Field", floatField.Type())
	// 根据名字查找字段
	fmt.Println("FieldByName(\"b\").Type", d.FieldByName("b").Type())
	// 根据索引查找值中, next字段的int字段的值
	fmt.Println("FieldByIndex([]int{4, 0}).Type()", d.FieldByIndex([]int{4, 0}).Type())

	x := reflect.ValueOf(&dummy{
		next: &dummy{},
	})

	fmt.Println("NumField-valueof", x.Elem().NumField())
	fmt.Println("NumField-valueof", x.Elem().Field(0).Type().Size())

	gg := reflect.TypeOf(&dummy{
		next: &dummy{},
	})
	fmt.Println("NumField-TypeOf", gg.Elem().NumField()) // 必须用结构体,所以要Elem下解引用,变成结构体,再查找里面的具体字段
	fmt.Println("NumField-TypeOf", gg.Elem().Field(0).Type.Size())
}

//演示  IsNil() 常被用于判断指针是否为空；IsValid() 常被用于判定返回值是否有效。
func NilAndValild() {
	fmt.Println("---------NilAndValild----------")
	// *int的空指针
	var a *int
	fmt.Println("var a *int:", reflect.ValueOf(a).IsNil())
	// nil值
	fmt.Println("nil:", reflect.ValueOf(nil).IsValid())
	// *int类型的空指针
	fmt.Println("(*int)(nil):", reflect.ValueOf((*int)(nil)).Elem().IsValid())
	// 实例化一个结构体
	s := struct{}{}
	// 尝试从结构体中查找一个不存在的字段
	fmt.Println("不存在的结构体成员:", reflect.ValueOf(s).FieldByName("").IsValid())
	// 尝试从结构体中查找一个不存在的方法
	fmt.Println("不存在的结构体方法:", reflect.ValueOf(s).MethodByName("").IsValid())
	// 实例化一个map
	m := map[int]int{}
	// 尝试从map中查找一个不存在的键
	fmt.Println("不存在的键：", reflect.ValueOf(m).MapIndex(reflect.ValueOf(3)).IsValid())
}

//可寻址的变量才能设置值
func setValue() {
	fmt.Println("---------setValue----------")
	type dog struct {
		LegCount int // 必须为导出
	}
	// 获取dog实例地址的反射值对象
	valueOfDog := reflect.ValueOf(&dog{})
	// 取出dog实例地址的元素
	valueOfDog = valueOfDog.Elem()
	// 获取legCount字段的值
	vLegCount := valueOfDog.FieldByName("LegCount")
	// 尝试设置legCount的值
	vLegCount.SetInt(4)
	fmt.Println(vLegCount.Int())

	// 切片 就是指针,可以直接用
	var a []byte = []byte{0x1, 0x2}
	fmt.Println(a)
	d := reflect.ValueOf(a)
	fmt.Println("d.Len()", d.Len())

	fmt.Println(d.CanSet())
	d.Index(0).SetUint(0x03)
	fmt.Println(a)

	fmt.Println("-------------", d.Type().Elem())

	//切片长度只用用ValueOf获取
	// fmt.Println("-------")
	// s2 := reflect.TypeOf(a)
	// fmt.Println("s.Len()", s2.Len())

	var arr [2]byte = [...]byte{0x1, 0x2}

	s := reflect.TypeOf(arr)
	fmt.Println("s.Len()", s.Len())

	s1 := reflect.ValueOf(arr)
	fmt.Println("s1.Len()", s1.Len())

	fmt.Println("s1.Elem()", s1.Type().Elem())
	//fmt.Println("s1.Elem()", s1.Elem()) // 错误
}

func newInstance() {
	var a int
	// 取变量a的反射类型对象
	typeOfA := reflect.TypeOf(a)
	// 根据反射类型对象创建类型实例
	aIns := reflect.New(typeOfA)
	// 输出Value的类型和种类
	fmt.Println(aIns.Type(), aIns.Kind()) // *int ptr
}

// 普通函数
func add(a, b int) int {
	return a + b
}
func call() {
	// 将函数包装为反射值对象
	funcValue := reflect.ValueOf(add)
	// 构造函数参数, 传入两个整型值
	paramList := []reflect.Value{reflect.ValueOf(10), reflect.ValueOf(20)}
	// 反射调用函数
	retList := funcValue.Call(paramList)
	// 获取第一个返回值, 取整数值
	fmt.Println(retList[0].Int())
}

func getPtrSize() {
	x := 3
	p := reflect.ValueOf(&x)
	fmt.Println("getPtrSize", p.Type().Size())
	fmt.Println("getPtrSize", unsafe.Sizeof(x))
}
