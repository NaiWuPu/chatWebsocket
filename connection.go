package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
)

// 抽象出需要的数据结构

//ws 连接器 数据 管道
type connection struct {
	ws   *websocket.Conn // 连接器
	send chan []byte     // 管道
	data *Data           // 数据
}

// 抽象 ws 连接器 处理ws 的各种逻辑
type hub struct {
	connections map[*connection]bool // connections 注册了连接器
	broadcast   chan []byte          // 从连接器发送的信息
	register    chan *connection     // 从连接器注册请求
	unregister  chan *connection     // 销毁请求
}

// ws的写
func (c *connection) writer() {
	for message := range c.send {
		//数据写出
		_ = c.ws.WriteMessage(websocket.TextMessage, message)
	}
	_ = c.ws.Close()
}

var userList []string

// ws的读
func (c *connection) reader() {
	// 不断的去读websocket数据
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			// 读不进数据，用户移除
			h.unregister <- c
			break
		}
		// 读取数据
		_ = json.Unmarshal(message, c.data)
		// 根据data 的type 判断该做什么
		switch c.data.Type {
		case "login":
			// 弹出窗口， 输入用户名
			c.data.User = c.data.Content
			c.data.From = c.data.User
			// 登录后，将用户加入列表
			userList = append(userList, c.data.User)
			// 每个用户都加载用户列表
			c.data.UserList = userList
			// 数据序列化
			dadaB, _ := json.Marshal(c.data)
			h.broadcast <- dadaB
			// 普通状态
		case "user":
			c.data.Type = "user"
			dataB, _ := json.Marshal(c.data)
			h.broadcast <- dataB
		case "logout":
			c.data.Type = "logout"
			userList = remove(c.data.UserList, c.data.User)
			c.data.UserList = userList
			c.data.Content = c.data.User
			// 数据序列化
			dataB, _ := json.Marshal(c.data)
			h.broadcast <- dataB
			h.unregister <- c
		default:

		}
	}
}

// 删除用户切片中的数据
func remove(slice []string, user string) []string {
	// 严谨判断
	count := len(slice)
	if count == 0 {
		return slice
	}
	if count == 1 && slice[0] == user {
		return []string{}
	}
	// 定义新的返回切片
	var mySlice []string
	for i := range slice {
		if slice[i] == user && i == count {
			return slice[:count]
		} else if slice[i] == user {
			mySlice = append(slice[:i], slice[i+1:]...)
			break
		}
	}
	return mySlice
}

// 定义升级器 将http请求升级为ws请求
var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	}}

// ws的回调函数
func wsHandler(w http.ResponseWriter, r *http.Request) {
	// 获取 ws 对象
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	// 创建连接对象去做事情
	c := &connection{send: make(chan []byte, 128), ws: ws, data: &Data{}}
	// 在ws注册一下
	h.register <- c
	// ws将数据读写起来
	go c.writer()
	c.reader()
	defer func() {
		c.data.Type = "logout"
		userList = remove(c.data.UserList, c.data.User)
		c.data.UserList = userList
		c.data.Content = c.data.User
		// 数据序列化
		dataB, _ := json.Marshal(c.data)
		h.broadcast <- dataB
		h.unregister <- c
	}()
}
