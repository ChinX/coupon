package memory

import (
	"errors"
	"sync"

	"github.com/chinx/coupon/session"
)

type Store struct {
	lock  sync.RWMutex
	data  map[string]*session.Session
	interval int64
}

// NewStore creates and returns a session store.
func NewStore() *Store {
	return &Store{
		data:  make(map[string]*session.Session),
	}
}

func (s *Store) Set(key string, val *session.Session, expire int64) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.data[key] = val
	return nil
}

func (s *Store) Get(key string) (*session.Session, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	val, ok := s.data[key]
	if ok {
		return val, nil
	}
	return nil, errors.New("key " + key + " is not found")
}

func (s *Store) Expire(key string, expire int64) error {

}

func (s *Store)gc()  {
	for {
		select {
		case :
			
		}
	}
}
