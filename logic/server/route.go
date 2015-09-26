package server

import (
	"encoding/json"
	"errors"
)

var roter map[string]func(*Conn) = make(map[string]func(*Conn))

func registerRoter() {
	//用户绑定
	roter["99"] = bind
	//用户注册
	roter["100"] = register
}
func handler(c *Conn, data []byte) {
	err := json.Unmarshal(data, &c.Request)
	if err != nil {
		c.Out([]byte(err.Error()))
		return
	}
	h := roter[c.Head.Command]
	if h == nil {
		c.Out([]byte(errors.New("Does not exist command?").Error()))
		return
	}
	//加入拦截器 检查身份
	h(c)
	//选择执行函数
}
