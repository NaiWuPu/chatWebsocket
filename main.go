package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	// 创建路由
	router := mux.NewRouter()
	go h.run()
	// 知道ws回调函数
	router.HandleFunc("/ws", wsHandler)
	// 开启服务端监听
	if err := http.ListenAndServe("127.0.0.1:8080", router); err != nil {
		fmt.Println("err:", err)
	}
}
