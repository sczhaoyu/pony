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
func (c *Conn) NewResponse(data []byte) *common.Response {
	var w common.Response
	w.Head = new(common.ResponseHead)
	w.Head.SessionId = c.Request.Head.SessionId
	w.Head.Command = c.Request.Head.Command
	w.Body = data
	return &w
}
func (c *Conn) Out(data []byte) {
	var w Write
	w.Conn = c
	w.Body = util.GetJsonByteLen(c.NewResponse(data))
	c.Put(&w)
}
