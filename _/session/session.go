package session

import (
	"net"
)

type Session struct {
	net.Conn                         //会话链接
	SESSIONID string                 //会话ID由服务器生成
	Attr      map[string]interface{} //会话中存储的值
}

func (c *Session) GetAttr(k string) interface{} {
	return c.Attr[k]
}
func (c *Session) SetAttr(k string, v interface{}) {
	if c.Attr == nil {
		c.Attr = make(map[string]interface{})
	}
	c.Attr[k] = v
}
