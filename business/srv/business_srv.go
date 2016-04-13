package srv

import (
	"fmt"
	"github.com/sczhaoyu/pony/business/srv/customer"
	. "github.com/sczhaoyu/pony/server"
)

type BusinessServer struct {
	Srv *Server         //链接服务
	CS  *CustomerServer //客户端服务(链接管理端)
}

//启动客户端服务
func (c *BusinessServer) RunBusinessServer() {
	c.Srv = new(Server)
	c.Srv.Name = "业务服务(1)"
	c.Srv.Handle = Handle
	c.Srv.Port = 9088
	c.Srv.MaxClient = 1000
	c.Srv.Run()

}

//链接管理后台
func (c *BusinessServer) RunCustomerServer() {
	c.CS = new(CustomerServer)
	c.CS.Name = "后台服务(1)"
	c.CS.ServerAddr = "127.0.0.1:9789"
	c.CS.Handler = customer.Handle
	c.CS.FirstSend = func() {
		//发送自己的信息
		ret := &struct {
			Name string `json:"name"`
			Id   string `json:"id"`
		}{
			Name: c.Srv.Name,
			Id:   c.Srv.Id,
		}
		c.CS.WriteJson(ret, 100)
	}
	c.CS.Run()
}

func Handle(c *Conn) {

	//存储数据库
	//发送至业务逻辑服务器
	fmt.Println(c.Request)

}
