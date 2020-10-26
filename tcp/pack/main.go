package main

import (
	"fmt"
)

//普通的封包和解包
func pack() {
	fmt.Println("原始包")
	packData := NewPacket([]byte(`{"funcid":10,"userid":"12345678900987654321","notice":10.0,"name":2}`), []byte("success~~"))
	fmt.Println(packData)

	result := packData.Pack()          // 整个结构体转换为字节流
	result[len(result)-2] = 0x34       // 更改下字节流 相当于将2 改成 4
	newPack := packData.UnPack(result) // 将字节流转换为结构体
	fmt.Println("新包")
	fmt.Println(newPack)

}



func main() {
	pack()

}
