package server

import (
	"encoding/json"
)

var roter map[int]func(*Request) = make(map[int]func(*Request))

func registerRoter() {
	//用户注册
	roter[100] = register
}
func handler(data []byte) {
	var r Request
	err := json.Unmarshal(data, &r)
	if err != nil {
		//输出错误到客户端
		return
	}
	h := roter[r.Head.FaceCode]
	if h == nil {
		//函数不存在 输出错误到客户端
		return
	}
	//加入拦截器 检查身份
	h(&r)
	//选择执行函数
}
