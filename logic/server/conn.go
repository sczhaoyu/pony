package server

import (
	"github.com/sczhaoyu/pony/util"
	"net"
)

type Conn struct {
	net.Conn
	Request
	*Server
}

func NewConn(conn net.Conn, s *Server) *Conn {
	var c Conn
	c.Conn = conn
	c.Server = s
	return &c
}

//生成响应
func (c *Conn) NewResponse(d interface{}) *Response {
	var w Response
	w.Head = new(ResponseHead)
	w.Head.Addr = c.RemoteAddr().String()
	w.Head.Uuid = util.GetUUID()
	w.Head.UserId = c.Head.UserId
	w.Head.UserAddr = c.Head.UserAddr
	switch err := d.(type) {
	case int:
		w.Head.Msg = ErrMsg[d.(int)]
		w.Head.State = d.(int)
	case error:
		w.Head.Msg = err.Error()
		w.Head.State = -1
	default:
		w.Head.State = 0
		w.Body = d
	}
	return &w
}
func (c *Conn) Write(data interface{}) {
	c.Put(c.NewResponse(data))
}
