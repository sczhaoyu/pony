package common

type SessionManager struct {
	Session map[string]*UsersSession
}

func (s *SessionManager) Init() {
	s.Session = make(map[string]*UsersSession)
}
