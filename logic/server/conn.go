package server

import (
	"github.com/sczhaoyu/pony/common"
	"github.com/sczhaoyu/pony/util"
	"net"
)

type Conn struct {
	net.Conn
	common.Request
	*Server
}

func NewConn(conn net.Conn, s *Server) *Conn {
	var c Conn
	c.Conn = conn
	c.Server = s
	return &c
}

//生成响应
func (c *Conn) NewResponse(d interface{}) *common.Response {
	var w common.Response
	w.Head = new(common.ResponseHead)
	w.Head.UserId = c.Head.UserId
	w.Head.SessionId = c.Request.Head.SessionId
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
func (c *Conn) Out(data interface{}) {
	var w Write
	w.Conn = c
	w.Body = util.GetJsonByteLen(c.NewResponse(data))
	c.Put(&w)
}
