package common

import (
	"github.com/sczhaoyu/pony/util"
	"net"
	"time"
)

type Conn struct {
	net.Conn
	UUID  string
	Ctime int64
}

func NewConn(conn net.Conn) *Conn {
	var c Conn
	c.Conn = conn
	c.Ctime = time.Now().Unix()
	c.UUID = util.GetUUID()
	return &c
}
