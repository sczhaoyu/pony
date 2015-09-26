package admin_server

import (
	"net"
)

type ClientSer struct {
	net.Conn
	Addr       string //服务器IP
	ServerType string //服务器类型

}
