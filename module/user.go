package module

import (
	"context"
	"crypto/sha1"
	"crypto/sha512"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/chinx/coupon/model"
	"github.com/go-session/session"
)

const (
	StatusLogout int = iota
	StatusBinding
	StatusLogin

	PermissionUser  int = 2
	PermissionAdmin int = 7

	idKey         = "openid"
	sessionKey    = "session_key"
	permissionKey = "permission"
)

type AdminLogin struct {
	User   string `json:"user"`
	Passwd string `json:"passwd"`
}

type Session struct {
	store session.Store
	w     http.ResponseWriter
	r     *http.Request
}

func NewSession(w http.ResponseWriter, r *http.Request) (*Session, error) {
	store, err := session.Start(context.Background(), w, r)
	if err != nil {
		return nil, err
	}
	return &Session{store: store, w: w, r: r}, nil
}

func (s *Session) IsAdmin() bool {
	return s.IsPermission(PermissionAdmin)
}

func (s *Session) IsUser() bool {
	return s.IsPermission(PermissionUser)
}

func (s *Session) IsPermission(permission int) bool {
	key, ok := s.store.Get(permissionKey)
	if !ok {
		return false
	}
	return int(key.(float64)) >= permission
}

func (s *Session) ValidateSignature(sign string, raw string) (int, error) {
	key, ok := s.store.Get(sessionKey)
	if !ok {
		return s.Destroy(), errors.New("获取登录信息失败")
	}

	if sign != signature(raw+key.(string)) {
		return StatusLogin, errors.New("请求参数不合法")
	}

	return StatusLogin, nil
}

func (s *Session) Binding(data *Binding) error {
	user := &model.User{}
	err := json.Unmarshal([]byte(data.RawData), user)
	if err != nil {
		return err
	}
	openid, _ := s.store.Get(idKey)
	condition := &model.User{ID: openid.(string)}
	user.ID = condition.ID

	has := model.Get(condition)
	if !has {
		ok := model.Insert(user)
		if !ok {
			return errors.New("绑定用户信息失败")
		}
	} else {
		if condition.AvatarURL != user.AvatarURL ||
			condition.City != user.City ||
			condition.Province != user.Province ||
			condition.Country != user.Country ||
			condition.Gender != user.Gender ||
			condition.Language != user.Language ||
			condition.Nickname != user.Nickname {
			ok := model.Update(user)
			if !ok {
				return errors.New("更新用户信息失败")
			}
		}
	}
	return nil
}

func (s *Session) SetUserSession(wxData *WXSession) int {
	s.store.Set(idKey, wxData.OpenID)
	s.store.Set(sessionKey, wxData.SessionKey)
	s.store.Set(permissionKey, PermissionUser)
	err := s.store.Save()
	if err != nil {
		return s.Destroy()
	}

	user := &model.User{ID: wxData.OpenID}
	if ok := model.Get(user); !ok {
		return StatusBinding
	}
	return StatusLogin
}

func (s *Session) SetAdminSession(data *AdminLogin) int {
	admin := &model.Admin{User: data.User}
	if ok := model.Get(admin); !ok {
		return s.Destroy()
	}

	sha := sha512.New()
	sha.Write([]byte(admin.User + admin.Salt))
	if fmt.Sprintf("%x", sha.Sum(nil)) != admin.Passwd {
		return s.Destroy()
	}

	s.store.Set(idKey, data.User)
	s.store.Set(permissionKey, PermissionAdmin)
	err := s.store.Save()
	if err != nil {
		return s.Destroy()
	}
	return StatusLogin
}

func (s *Session) Destroy() int {
	session.Destroy(context.Background(), s.w, s.r)
	return StatusLogout
}

func (s *Session) UserID() (string, int) {
	userID, ok := s.store.Get(idKey)
	if !ok {
		return "", s.Destroy()
	}

	return userID.(string), StatusLogin
}

func signature(key string) string {
	h := sha1.New()
	io.WriteString(h, key)
	return fmt.Sprintf("%x", h.Sum(nil))
}
