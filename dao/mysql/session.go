package mysql

import "github.com/go-xorm/xorm"

type Session struct {
	session *xorm.Session
}

func NewSession() *Session {
	return &Session{session: engine.NewSession()}
}

func (s *Session) Begin() error {
	return s.session.Begin()
}

func (s *Session) Rollback() error {
	return s.session.Rollback()
}

func (s *Session) Commit() error {
	return s.session.Commit()
}

func (s *Session) Close() {
	s.session.Close()
}

func (s *Session) Insert(tb tabler) (err error) {
	return insert(s.session, tb)
}

func (s *Session) Update(tb tabler, condition string, args ...interface{}) (err error) {
	return update(s.session, tb, condition, args...)
}

func (s *Session) Find(tb tabler, condition string, args ...interface{}) (list []interface{}, err error) {
	return find(s.session, tb, condition, args...)
}

func (s *Session) In(tb tabler, field string, conditions []interface{}) []interface{} {
	return in(s.session, tb, field, conditions)
}

func (s *Session) List(tb tabler, offset, count int, condition string, args ...interface{}) (int64, []interface{}) {
	return list(s.session, tb, offset, count, condition, args...)
}

func (s *Session) GetLast(tb tabler, order, condition string, args ...interface{}) (err error) {
	return getLast(s.session, tb, order, condition, args...)
}

func (s *Session) Get(tb tabler, condition string, args ...interface{}) (err error) {
	return get(s.session, tb, condition, args...)
}

func (s *Session) Exist(tb tabler, condition string, args ...interface{}) bool {
	return exist(s.session, tb, condition, args...)
}

func (s *Session) Delete(tb tabler, condition string, args ...interface{}) (err error) {
	return delete(s.session, tb, condition, args...)
}
