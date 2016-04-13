package srv

import (
	. "github.com/sczhaoyu/pony/server"
)

//路由处理
var router map[int]func(c *Conn) = make(map[int]func(c *Conn))

func init() {
	//逻辑服务器登录
	router[100] = loginBusiness
	router[520] = herbat
	router[200] = ok
}
