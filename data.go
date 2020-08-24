package main

type Data struct {
	Ip   string `json:"ip"`
	Type string `json:"type"` // 标识信息的类型
	// login 登录
	// handshake
	// system
	// logout
	// user
	From    string `json:"from"`    // 代表哪个用户
	Content string `json:"content"` // 传送内容
	// 用户名
	User string `json:"user"`
	// 用户列表
	UserList []string `json:"user_list"`
}
