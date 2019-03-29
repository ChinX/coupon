package session

import (
	"github.com/chinx/coupon/api"
	"github.com/chinx/coupon/module"
	"github.com/chinx/coupon/utils/httputils"
	"log"
	"net/http"
)

type Session struct {
	exception []string
}

func NewSession(exception ...string) *Session {
	return &Session{exception:exception}
}

func (s *Session) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	switch r.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		for _, val := range s.exception{
			if val == r.URL.Path {
				next(rw, r)
				return
			}
		}
		checkSession(rw, r, module.PermissionUser)
	default:
	}
	next(rw, r)
}

func checkSession(w http.ResponseWriter, r *http.Request, permission int) {
	result := &api.ReplyResult{Status: module.StatusLogout, ErrorCode: module.NoError, Surplus: module.LimitCount}
	userData, err := module.NewSession(w, r)
	if err != nil {
		log.Println("checkLogin", err)
		result.Message = "获取登录信息失败"
		httputils.Reply(w, http.StatusUnauthorized, result, err)
		return
	}

	userData.ShowALL()
	if !userData.IsPermission(permission) {
		result.Message = "无权限"
		httputils.Reply(w, http.StatusForbidden, result, err)
		return
	}

	result.UserID, result.Status = userData.UserID()
	if result.Status == module.StatusLogout {
		result.Message = "登录状态已失效，请重新登录"
		httputils.Reply(w, http.StatusUnauthorized, result, nil)
		return
	}
}