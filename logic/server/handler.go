package server

import (
	"encoding/json"
	"github.com/sczhaoyu/pony/common"
)

var sysRoter map[string]func(c *Conn) = make(map[string]func(c *Conn))

func steupSysRoter() {
	//注册链接池的所属客户端服务器
	sysRoter[common.LOGICCLIENT] = registerSession
}
func handler(c *Conn, data []byte) {
	err := json.Unmarshal(data, &c.Request)
	if err != nil {
		c.Out([]byte(err.Error()))
		return
	}
	if sysRoter[c.Head.Command] != nil {
		sysRoter[c.Head.Command](c)
		//系统通信直接处理不放行
		return
	}
	b := beforeInterceptor
	for i := 0; i < len(b); i++ {
		if b[i](c) == false {
			return
		}
	}
	if ReadFunc != nil {
		ReadFunc(c)
	}
	a := afterInterceptor
	for i := 0; i < len(a); i++ {
		if a[i](c) == false {
			return
		}
	}
}
func registerSession(c *Conn) {
	addr := string(c.Body)
	//加入正式会话session
	c.Session.SetSession(c, addr)
}
