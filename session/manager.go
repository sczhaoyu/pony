package session

import (
	"github.com/sczhaoyu/pony/util"
	"net"
	"sync"
)

type Manager struct {
	Session map[string]*Session
	sm      sync.Mutex
}

func (m *Manager) Init() {
	m.Session = make(map[string]*Session)
}
func (m *Manager) GetSession(key string) *Session {
	return m.Session[key]
}
func (m *Manager) SetSession(conn net.Conn) *Session {
	m.sm.Lock()
	var s Session
	s.SESSIONID = util.GetUUID()
	s.Conn = conn
	m.Session[s.SESSIONID] = &s
	defer m.sm.Unlock()
	return &s
}
func (m *Manager) Remove(s *Session) {
	m.sm.Lock()
	s.Close()
	delete(m.Session, s.SESSIONID)
	m.sm.Unlock()
}
