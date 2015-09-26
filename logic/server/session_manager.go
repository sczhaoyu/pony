package server

import (
	"log"
	"net"
	"sync"
)

//根据sessionId拿到输出管道
//全局管道获取
//部分用户管道获取
type SessionManager struct {
	SCName    map[string]map[string]net.Conn //记录和Addr的关系
	SessionId map[string]string              //记录和SESSIONID的关系
	mutex     sync.Mutex                     //会话操作锁
}

func (s *SessionManager) Init() {
	s.SCName = make(map[string]map[string]net.Conn)
	s.SessionId = make(map[string]string)
}

//根据sessionId获取session
func (s *SessionManager) GetSessionId(sid string) net.Conn {
	id := s.SessionId[sid]
	if id == "" {
		return nil
	}
	sc := s.SCName[id]
	for _, v := range sc {
		return v
	}
	return nil
}

//添加session
func (s *SessionManager) SetSession(c net.Conn, scName string) {
	s.mutex.Lock()
	if s.SCName[scName] == nil {
		s.SCName[scName] = make(map[string]net.Conn)
	}
	sc := s.SCName[scName]
	switch conn := c.(type) {
	case *Conn:
		//查看临时会话中是否存在如果存在移除
		if sc[conn.RemoteAddr().String()] != nil {
			delete(sc, conn.RemoteAddr().String())
		}
		sc[conn.RemoteAddr().String()] = conn
		s.SessionId[conn.Head.SessionId] = scName
	default:
		sc[c.RemoteAddr().String()] = c
	}
	defer s.mutex.Unlock()
	log.Println(s.SCName)
}

func (s *SessionManager) RemoveAddrSession(addr string) {
	s.mutex.Lock()
	for _, v := range s.SCName {
		for a, c := range v {
			if a == addr {
				//删除
				c.Close()
				delete(v, addr)
			}
		}
	}
	defer s.mutex.Unlock()
	log.Println(s.SCName)
}
