package main

import "encoding/json"

// 连接器对象初始化
var h = hub{
	connections: make(map[*connection]bool), // connections 注册了连接器
	broadcast:   make(chan []byte),          // 从连接器发送的信息
	register:    make(chan *connection),     // 从连接器注册请求
	unregister:  make(chan *connection),     // 销毁请
}

// 处理ws 的逻辑实现
func (h *hub) run() {
	// 监听数据管道，在后端不断处理管道数据
	for {
		// 根据不同的数据管道处理不同逻辑
		select {
		case c := <-h.register:
			// 标识注册了
			h.connections[c] = true
			c.data.Ip = c.ws.RemoteAddr().String()
			c.data.Type = "handshake"
			c.data.UserList = userList
			dataB, _ := json.Marshal(c.data)
			c.send <- dataB
		case c := <-h.unregister:
			// 判断map 里灿在处理的数据
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
		case data := <-h.broadcast:
			// 处理数据流转，主句同步到所有用户
			// c 是具体的每一个连接
			for c := range h.connections {
				select {
				case c.send <- data:
				default:
					// 防止死循环
					delete(h.connections, c)
					close(c.send)
				}
			}
		}
	}
}
