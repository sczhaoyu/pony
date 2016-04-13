package srv

import (
	"errors"
	. "github.com/sczhaoyu/pony/server"
	"log"
)

type BusinessSrv struct {
	Name  string //业务服务器名称
	Id    string //业务服务器标示
	At    int64  //链接时间
	*Conn        //链接
}

//全部业务处理服务器
var BS map[string]*BusinessSrv = make(map[string]*BusinessSrv, 0)

//加入业务服务器列表
func AddBusinessSrv(conn *Conn, name, id string) {
	var b BusinessSrv
	b.Conn = conn
	b.Name = name
	b.Id = id
	BS[id] = &b
}

//逻辑服务器登录
func loginBusiness(c *Conn) {
	name := c.Value("name").ToString()
	id := c.Value("id").ToString()
	if name == "" || id == "" {
		c.WriteJson(errors.New("name or id not null!"))
		return
	}

	AddBusinessSrv(c, name, id)
	c.WriteJson(nil)
}

//心跳
func herbat(c *Conn) {
	//保持心跳,刷新session的时间
	c.Server.MemProvider.SessionUpdate(c.Session.SessionID())
}

//数据包收包应答
func ok(c *Conn) {
	//取出数据包的ID
	responsId := string(c.Body)
	//让数据包重发管理器删除
	c.Server.DPM.Receive(responsId)
	//不需要回复
	log.Println("received response id:", responsId)
}
