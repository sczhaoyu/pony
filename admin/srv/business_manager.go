package srv

import (
	"errors"
	"fmt"
	. "github.com/sczhaoyu/pony/config"
	. "github.com/sczhaoyu/pony/server"
)

type Srv struct {
	Name    string //业务服务器名称
	Id      string //业务服务器标示
	At      int64  //链接时间
	Addr    string //服务器的地址
	SrvType int    //服务器类型
}

//创建一个服务器基本信息
func NewSrv(name, id, addr string, srvType int) *Srv {
	var b Srv
	b.Name = name
	b.Id = id
	b.Addr = addr
	b.SrvType = srvType
	return &b
}

//逻辑服务器登录
func loginServer(c *Conn) {

	name := c.Value("name").ToString()
	id := c.Value("id").ToString()
	addr := c.Value("addr").ToString()
	srvType := c.Value("srvType").ToInt()
	if name == "" || id == "" {
		c.WriteJson(errors.New("name or id not null!"))
		return
	}
	if srvType != SERVER_TYPE_BS && SERVER_TYPE_CS != srvType {
		c.WriteJson(errors.New("server type not found!"))
		return
	}

	//刷新session
	c.Session.Set(SERVER_TYPE_KEY, NewSrv(name, id, addr, srvType))
	c.WriteJson(nil)
	//通知客户端服务器,从session里面通知
	if srvType == SERVER_TYPE_BS {
		noticeSrvCS([]string{addr}, SERVER_TYPE_CS, 101)
	}
}

//通知客户端服务器组
func noticeSrvCS(b interface{}, srvType, faceCode int) {
	ret := Admin.Server.MemProvider.SessionStoreAll()
	for i := 0; i < len(ret); i++ {
		v := ret[i].Get(SERVER_TYPE_KEY)
		if v != nil {
			if v.(*Srv).SrvType == srvType {
				ret[i].Notice(b, faceCode)
			}
		}
	}

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
	srv := c.Session.Get(SERVER_TYPE_KEY).(*Srv)
	ret := fmt.Sprintf("%s:%s", srv.Name, responsId)
	fmt.Println(ret)
}

//获取在线逻辑服务器列表
func getBusinessList(c *Conn) {
	ret := Admin.Server.MemProvider.SessionStoreAll()
	adrs := []string{}
	for i := 0; i < len(ret); i++ {
		v := ret[i].Get(SERVER_TYPE_KEY)
		if v != nil {
			if v.(*Srv).SrvType == SERVER_TYPE_BS {
				adrs = append(adrs, v.(*Srv).Addr)
			}
		}
	}

	c.WriteJson(adrs)
}
