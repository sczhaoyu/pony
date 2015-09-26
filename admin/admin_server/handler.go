package admin_server

import (
	"github.com/sczhaoyu/pony/common"
	"net"
)

func handler(a *AdminServer, conn net.Conn, data []byte) {
	//业务处理
	var r common.Request
	r.Unmarshal(data)
	switch r.Head.Command {
	case common.LS:
		a.AddSession(conn, common.LS)
		//通知前端服务器有新的逻辑服务器加入
		sp := common.AuthResponse(common.ADDLS, r.Body)
		a.SendNotice(common.CS, sp.GetJson())
	case common.CS:
		a.AddSession(conn, common.CS)

	}
}
