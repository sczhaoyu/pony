package session

import (
	"github.com/sczhaoyu/pony/util"
	"net"
	"sync"
)

type Manager struct {
	session map[string]*Session
	sm      sync.Mutex
}

func (m *Manager) Init() {
	m.session = make(map[string]*Session)
}
func (m *Manager) GetSession(key string) *Session {
	return m.session[key]
}
func (m *Manager) SetSession(conn net.Conn) *Session {
	m.sm.Lock()
	var s Session
	s.SESSIONID = util.GetUUID()
	s.Conn = conn
	m.session[s.SESSIONID] = &s
	defer m.sm.Unlock()
	return &s
}
func (m *Manager) Remove(s *Session) {
	m.sm.Lock()
	delete(m.session, s.SESSIONID)
	m.sm.Unlock()
}
