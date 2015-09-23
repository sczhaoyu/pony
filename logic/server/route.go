package server

import (
	"encoding/json"
	"errors"
)

var roter map[int]func(*Conn) = make(map[int]func(*Conn))

func registerRoter() {
	//用户注册
	roter[100] = register
}
func handler(c *Conn, data []byte) {
	err := json.Unmarshal(data, &c.Request)
	if err != nil {
		c.Out(err)
		return
	}
	h := roter[c.Head.FaceCode]
	if h == nil {
		c.Out(errors.New("Does not exist faceCode?"))
		return
	}
	//加入拦截器 检查身份
	h(c)
	//选择执行函数
}
