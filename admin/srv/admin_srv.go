package srv

import (
	. "github.com/sczhaoyu/pony/server"
)

type AdminServer struct {
	Server //嵌入服务器

}

var Admin *AdminServer

func Run() {
	Admin = new(AdminServer)
	Admin.Name = "后台服务器"
	Admin.Port = 9789
	Admin.Handle = RouterHandle
	Admin.Run()
}

//路由处理
func RouterHandle(c *Conn) {
	//根据路由命令转发响应的功能
	fun, ok := router[c.Request.Header.FaceCode]
	if !ok {
		c.WriteJson("not found faceCode")
		return
	}
	fun(c)

}
