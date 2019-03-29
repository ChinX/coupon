package module

import (
	"context"
	"crypto/sha1"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/chinx/coupon/model"

	"github.com/chinx/coupon/dao/mysql"
	"github.com/go-session/session"
)

const (
	StatusLogout int = iota
	StatusBinding
	StatusLogin

	PermissionUser  int = 2
	PermissionAdmin int = 7

	userIdKey     = "user_id"
	openIdKey     = "open_id"
	sessionKey    = "session_key"
	permissionKey = "permission"
)

type AdminLogin struct {
	User     string `json:"user"`
	Password string `json:"password"`
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
		return StatusBinding, errors.New("请求参数不合法")
	}

	return StatusBinding, nil
}

func (s *Session) Binding(data *Binding) (*model.User, error) {
	user := &model.User{}
	err := json.Unmarshal([]byte(data.RawData), user)
	if err != nil {
		return nil, err
	}
	log.Println(*user)
	userID, _ := s.store.Get(userIdKey)
	user.ID = int64(userID.(float64))
	user.Binding = 1

	condition := &model.User{}
	if err := mysql.Get(condition, "id=?", user.ID); err != nil {
		log.Println(err)
		return nil, errors.New("绑定用户信息失败")
	} else {
		if condition.AvatarURL != user.AvatarURL ||
			condition.City != user.City ||
			condition.Province != user.Province ||
			condition.Country != user.Country ||
			condition.Gender != user.Gender ||
			condition.Language != user.Language ||
			condition.Nickname != user.Nickname ||
			condition.Binding != user.Binding{
			err := mysql.Update(user, "id=?", user.ID)
			if err != nil && err != mysql.NoRecords {
				return nil, errors.New("更新用户信息失败")
			}
		}
	}
	return user, nil
}

func (s *Session) SetUserSession(wxData *WXSession) int {
	s.store.Set(userIdKey, wxData.ID)
	s.store.Set(openIdKey, wxData.OpenID)
	s.store.Set(sessionKey, wxData.SessionKey)
	s.store.Set(permissionKey, PermissionUser)
	err := s.store.Save()
	if err != nil {
		return s.Destroy()
	}

	user := &model.User{}
	if err := mysql.Get(user, "id=?", wxData.ID); err != nil || user.Binding == 0{
		return StatusBinding
	}
	return StatusLogin
}

func (s *Session) SetAdminSession(data *AdminLogin) int {
	admin := &model.Admin{User: data.User}
	if err := mysql.Get(admin, "user=?", data.User); err != nil {
		return s.Destroy()
	}

	sha := sha512.New()
	sha.Write([]byte(data.Password + admin.Salt))
	if hex.EncodeToString(sha.Sum(nil)) != admin.Password {
		return s.Destroy()
	}

	s.store.Set(userIdKey, data.User)
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

func (s *Session) Refresh(w http.ResponseWriter, r *http.Request) {
	w.Header().Del("Set-Cookie")
	r.Header.Del("Cookie")
	s.w = w
	s.r = r
	s.store, _ = session.Refresh(context.Background(), s.w, s.r)
}

func (s *Session) OpenID() (string, int) {
	userID, ok := s.store.Get(openIdKey)
	if !ok {
		return "", s.Destroy()
	}

	return userID.(string), StatusLogin
}

func (s *Session) UserID() (int64, int) {
	userID, ok := s.store.Get(userIdKey)
	if !ok {
		return 0, s.Destroy()
	}

	return int64(userID.(float64)), StatusLogin
}

func (s *Session) ShowALL() {
	log.Println(userIdKey)
	log.Println(s.store.Get(userIdKey))
	log.Println(sessionKey)
	log.Println(s.store.Get(sessionKey))
	log.Println(permissionKey)
	log.Println(s.store.Get(permissionKey))
}

func signature(key string) string {
	sha := sha1.New()
	io.WriteString(sha, key)
	return hex.EncodeToString(sha.Sum(nil))
}
