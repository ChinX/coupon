package session

import (
	"sync"
	"time"
)

type Provider interface {
	Set(string, *Session, int64) error
	Get(key string) (*Session, error)
	Delete(key string) error
	Expire(key string, expire int64) error
}

type Session struct {
	p      Provider
	sid    string
	expire int64
	Start  time.Time
	lock   sync.RWMutex
	data   map[interface{}]interface{}
}

// NewSession creates and returns a session store.
func NewSession(provider Provider, sid string, expire int64, data map[interface{}]interface{}) *Session {
	return &Session{
		p:      provider,
		expire: expire,
		sid:    sid,
		data:   data,
		Start:  time.Now(),
	}
}

// Set sets value to given key in session.
func (s *Session) Set(key, val interface{}) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.data[key] = val
	return nil
}

// Get gets value by given key in session.
func (s *Session) Get(key interface{}) interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.data[key]
}

// Delete delete a key from session.
func (s *Session) Delete(key interface{}) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.data, key)
	return nil
}

// ID returns current session ID.
func (s *Session) ID() string {
	return s.sid
}

// Release releases resource and save data to provider.
func (s *Session) Release() error {
	s.Start = time.Now()
	return s.p.Set(s.sid, s, s.expire)
}

func (s *Session) Expire() error {
	return s.p.Expire(s.sid, s.expire)
}

// Flush deletes all session data.
func (s *Session) Flush() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.data = make(map[interface{}]interface{})
	return nil
}
