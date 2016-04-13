package srv

import (
	"fmt"
	. "github.com/sczhaoyu/pony/server"
	"github.com/sczhaoyu/pony/util"
)

type ClientServer struct {
	Srv *Server         //链接服务
	CS  *CustomerServer //客户端服务
}

//启动客户端服务
func (c *ClientServer) RunSrv() {
	c.Srv = new(Server)
	c.Srv.Name = "客户端(1)"
	c.Srv.Handle = c.Handle
	c.Srv.MaxClient = 1
	c.Srv.Run()

}

//启动客户端服务
func (c *ClientServer) RunCustomerServer() {
	c.CS = new(CustomerServer)
	c.CS.Name = "后台服务(1)"
	c.CS.ServerAddr = "127.0.0.1:9789"
	c.CS.Handler = HandleCustomer
	c.CS.Run()
}
func HandleCustomer(conn *CustomerServer, rsp *Respon) {
	fmt.Println(rsp)

}
func (c *ClientServer) Handle(conn *Conn) {
	c.CS.DataChan <- util.ByteLen(conn.Request.Marshal())

}
