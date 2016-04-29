package srv

import (
	. "github.com/sczhaoyu/pony/server"
)

//路由处理
var router map[int]func(c *Conn) = make(map[int]func(c *Conn))

func init() {
	//服务器登录
	router[100] = loginServer
	//获取业务服务器列表
	router[101] = getBusinessList
	//确认收到数据包
	router[200] = ok
	//心跳包
	router[520] = herbat

}
