//服务器端session管理
package server

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

type SessionStore struct {
	sid          string                      //session id
	timeAccessed time.Time                   //last access time
	value        map[interface{}]interface{} //session store
	lock         sync.RWMutex
	*Conn        //链接信息
}

//设置会话信息
func (st *SessionStore) Set(key, value interface{}) error {
	st.lock.Lock()
	defer st.lock.Unlock()
	st.value[key] = value
	return nil
}

// 获取会话信息
func (st *SessionStore) Get(key interface{}) interface{} {
	st.lock.RLock()
	defer st.lock.RUnlock()
	if v, ok := st.value[key]; ok {
		return v
	} else {
		return nil
	}
}

// delete in memory session by key
func (st *SessionStore) Delete(key interface{}) error {
	st.lock.Lock()
	defer st.lock.Unlock()
	delete(st.value, key)
	return nil
}

// 请空话信息
func (st *SessionStore) Flush() error {
	st.lock.Lock()
	defer st.lock.Unlock()
	st.value = make(map[interface{}]interface{})
	return nil
}

// get 自己的sessionID
func (st *SessionStore) SessionID() string {
	return st.sid
}

type MemProvider struct {
	lock        sync.RWMutex             // locker
	sessions    map[string]*list.Element // map in memory
	list        *list.List               // for gc
	maxlifetime int64
}

// init  session
func (pder *MemProvider) SessionInit(maxlifetime int64) {
	pder.maxlifetime = maxlifetime

}

//根据sid获取session
func (pder *MemProvider) SessionRead(sid string) (*SessionStore, error) {
	pder.lock.RLock()
	if element, ok := pder.sessions[sid]; ok {
		pder.lock.RUnlock()
		return element.Value.(*SessionStore), nil
	} else {
		pder.lock.RUnlock()
		pder.lock.Lock()
		newsess := &SessionStore{sid: sid, timeAccessed: time.Now(), value: make(map[interface{}]interface{})}
		element := pder.list.PushBack(newsess)
		pder.sessions[sid] = element
		pder.lock.Unlock()
		return newsess, nil
	}
}

// 检测 sid是否存在
func (pder *MemProvider) SessionExist(sid string) bool {
	pder.lock.RLock()
	defer pder.lock.RUnlock()
	if _, ok := pder.sessions[sid]; ok {
		return true
	} else {
		return false
	}
}

// 销毁session
func (pder *MemProvider) SessionDestroy(sid string) error {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	if element, ok := pder.sessions[sid]; ok {
		session := element.Value.(*SessionStore)
		if session.Conn != nil {
			session.Close()
		}
		delete(pder.sessions, sid)
		pder.list.Remove(element)
		return nil
	}
	return nil
}

//session的GC
func (pder *MemProvider) SessionGC() {
	pder.lock.RLock()
	for {
		element := pder.list.Back()
		if element == nil {
			break
		}
		session := element.Value.(*SessionStore)
		if (session.timeAccessed.Unix() + pder.maxlifetime) < time.Now().Unix() {
			pder.lock.RUnlock()
			pder.lock.Lock()
			//告诉客户端登录超时
			session.WriteJson(errors.New("login time out"))
			//关闭客户端连接
			session.Close()
			pder.list.Remove(element)
			delete(pder.sessions, element.Value.(*SessionStore).sid)
			pder.lock.Unlock()
			pder.lock.RLock()
		} else {
			break
		}
	}
	pder.lock.RUnlock()
}

//获取session的数量
func (pder *MemProvider) SessionAll() int {
	return pder.list.Len()
}

//更新session的时间
func (pder *MemProvider) SessionUpdate(sid string) error {
	pder.lock.Lock()
	defer pder.lock.Unlock()
	if element, ok := pder.sessions[sid]; ok {
		element.Value.(*SessionStore).timeAccessed = time.Now()
		pder.list.MoveToFront(element)
		return nil
	}
	return nil
}

//一直回收
func (m *MemProvider) GC() {
	m.SessionGC()
	time.AfterFunc(time.Duration(m.maxlifetime)*time.Second, func() {
		m.GC()
	})
}
func NewSessionStore(maxlifetime int64) *MemProvider {
	var m MemProvider
	m.SessionInit(maxlifetime)
	m.list = list.New()
	m.sessions = make(map[string]*list.Element)
	go m.GC()
	return &m
}
