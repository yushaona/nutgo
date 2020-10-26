package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yushaona/nutgo/websocket/chatroom/controller"
)

var (
	addr   = ":2022"
	banner = `
	
    ____                _____
   |     |    |    /\     |
   |     |____|   /  \    | 
   |     |    |  /----\   |
   |____ |    | /      \  |

Go语言编程之旅 —— 一起用Go做项目：ChatRoom，start on：%s
`
)

func main() {
	fmt.Printf(banner, addr)

	controller.RegisterHandle()
	log.Fatal(http.ListenAndServe(addr, nil))
}
