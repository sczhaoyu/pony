package admin_server

import (
	"github.com/sczhaoyu/pony/common"
	"log"
	"net"
)

var roter map[string]func(*AdminServer, net.Conn, *common.Request) = make(map[string]func(*AdminServer, net.Conn, *common.Request))

func sutepRoter() {
	roter[common.LS] = regLS
	roter[common.CS] = regCS
	roter[common.GETLS] = getLS
	roter[common.ADDLSCONN] = addLogicConn
	roter[common.DELLSCONN] = delLogicConn
}

//登记逻辑服务器
func regLS(a *AdminServer, conn net.Conn, r *common.Request) {
	a.AddSession(conn, common.LS, string(r.Body))
	//通知前端服务器有新的逻辑服务器加入
	sp := common.AuthResponse(common.ADDLS, string(r.Body))
	a.SendNotice(common.CS, sp.GetJson())
}

//登记客户端链接服务器
func regCS(a *AdminServer, conn net.Conn, r *common.Request) {
	a.AddSession(conn, common.CS, string(r.Body))
}

//获取逻辑服务器组
func getLS(a *AdminServer, conn net.Conn, r *common.Request) {
	sp := common.AuthResponse(common.GETLS, a.GetLS())
	conn.Write(sp.GetJson())
}

//添加逻辑服务器session
func addLogicConn(a *AdminServer, conn net.Conn, r *common.Request) {
	a.mutex.Lock()
	s := a.CS[conn.RemoteAddr().String()]
	if s != nil {
		s.ClientNum = s.ClientNum + 1
	}
	defer a.mutex.Unlock()

}

//删除逻辑服务器session
func delLogicConn(a *AdminServer, conn net.Conn, r *common.Request) {
	a.mutex.Lock()
	s := a.CS[conn.RemoteAddr().String()]
	if s != nil {
		s.ClientNum = s.ClientNum - 1
	}
	defer a.mutex.Unlock()
}
func handler(a *AdminServer, conn net.Conn, data []byte) {
	//业务处理
	var r common.Request
	r.Unmarshal(data)
	cmd := r.Head.Command
	if roter[cmd] != nil {
		roter[cmd](a, conn, &r)
	}

}
