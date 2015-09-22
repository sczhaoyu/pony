package common

import (
	"net"
)

//会话管理
type UsersSession struct {
	UserId     int64    //用户的ID
	ClientConn net.Conn //会话链接
	UserAddr   string   //用户的IP信息
	UUID       string   //服务器端生成的UUID
}
