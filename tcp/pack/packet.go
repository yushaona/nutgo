package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
)

//DataPacket 构造协议头结构体
type DataPacket struct {
	Length     uint32
	Order      uint64
	Cmd        uint16
	StaticCode uint16 //0-正常 1-出错 2-命令未完成 0x000 3-命令需要ack
	Encryption uint8  // 0-误操作 1-AES256 2-国产sm4
	Reserv     uint32 //随机值
	InfoLen    uint16 //Info的长度 --  只有一个总长度Length没用
	Info       []byte //
	DataPack   []byte //数据包的长度
}

//Pack 打包成字节流
func (t *DataPacket) Pack() []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, t.Length)
	binary.Write(buf, binary.BigEndian, t.Order)
	binary.Write(buf, binary.BigEndian, t.Cmd)
	binary.Write(buf, binary.BigEndian, t.StaticCode)
	binary.Write(buf, binary.BigEndian, t.Encryption)
	binary.Write(buf, binary.BigEndian, t.Reserv)
	binary.Write(buf, binary.BigEndian, t.InfoLen)
	binary.Write(buf, binary.BigEndian, t.Info)
	binary.Write(buf, binary.BigEndian, t.DataPack)
	return buf.Bytes()
}

//UnPack 解析数据包
func (t *DataPacket) UnPack(packet []byte) *DataPacket {
	c := t.clone()
	buf := bytes.NewBuffer(packet)
	binary.Read(buf, binary.BigEndian, &c.Length)
	binary.Read(buf, binary.BigEndian, &c.Order)
	binary.Read(buf, binary.BigEndian, &c.Cmd)
	binary.Read(buf, binary.BigEndian, &c.StaticCode)
	binary.Read(buf, binary.BigEndian, &c.Encryption)
	binary.Read(buf, binary.BigEndian, &c.Reserv)
	binary.Read(buf, binary.BigEndian, &c.InfoLen)
	c.Info = make([]byte, c.InfoLen)
	binary.Read(buf, binary.BigEndian, &c.Info)                           //Read读取的时候,是根据data参数的大小,决定一次读取多少自己,看源码可知
	c.DataPack = make([]byte, c.Length-(4+8+2+2+1+4+2+uint32(c.InfoLen))) //重新初始化
	binary.Read(buf, binary.BigEndian, &c.DataPack)
	return c
}

func (t *DataPacket) String() string {
	return fmt.Sprintf("%+v\n", *t)
}

func (t *DataPacket) clone() *DataPacket {
	clone := *t
	return &clone
}

//NewPacket 构造数据包
func NewPacket(data []byte, info []byte) *DataPacket {
	infoLen := len(info)
	var dataPacket = &DataPacket{
		Length:     0,
		Order:      256,
		Cmd:        15,
		StaticCode: 3,
		Encryption: 1,
		Reserv:     2,
		InfoLen:    uint16(infoLen),
		Info:       info,
		DataPack:   data,
	}
	t := reflect.ValueOf(dataPacket).Elem()
	length := getSize(t) //  整个结构体的大小(字节)
	// l := 4 + 8 + 2 + 2 + 1 + 4 + 2 + infoLen + len(data)
	// fmt.Println(l)
	fmt.Println("NewPacket length", length)
	dataPacket.Length = uint32(length)
	return dataPacket
}

func getSize(t reflect.Value) int {
	switch t.Kind() {
	case reflect.Slice:
		n := sizeof(t.Type().Elem())
		return n * t.Len() // 切片的长度,只能通过reflect.Value获取
	case reflect.Struct:
		sum := 0
		for i, n := 0, t.NumField(); i < n; i++ { // 提取每个成员
			sum += getSize(t.Field(i))
		}
		return sum
	default:
		return sizeof(t.Type())
	}
}

func sizeof(t reflect.Type) int {
	//基本类型
	switch t.Kind() {
	case reflect.Ptr:
		return 1
	case reflect.Bool,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return int(t.Size())
	case reflect.Array:
		return sizeof(t.Elem()) * t.Len() // 数组的长度是类型的一部分,所以可以通过feflect.Type获取到
	case reflect.Struct:
		sum := 0
		for i, n := 0, t.NumField(); i < n; i++ {
			sum += sizeof(t.Field(i).Type)
		}
		return sum
	}
	return 0
}
