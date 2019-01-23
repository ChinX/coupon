package handler

import (
	"log"
	"net/http"

	"github.com/chinx/coupon/api"
	"github.com/chinx/coupon/module"
)

func UserSession(w http.ResponseWriter, r *http.Request) {
	operation := "[UserSession]"
	log.Println()
	log.Println(operation, "request url:", r.URL.String())
	log.Println("Cookie:", r.Header.Get("Cookie"))

	switch r.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		checkSession(w, r, module.PermissionUser)
	}
}

func AdminSession(w http.ResponseWriter, r *http.Request) {
	checkSession(w, r, module.PermissionAdmin)
}

func checkSession(w http.ResponseWriter, r *http.Request, permission int) {
	result := &api.ReplyResult{Status: module.StatusLogout}
	userData, err := module.NewSession(w, r)
	if err != nil {
		log.Println("checkLogin", err)
		result.Message = "获取登录信息失败"
		reply(w, http.StatusUnauthorized, result, err)
		return
	}

	log.Println("checkLogin")
	userData.ShowALL()
	if !userData.IsPermission(permission) {
		result.Message = "无权限"
		reply(w, http.StatusForbidden, result, err)
		return
	}

	result.UserID, result.Status = userData.UserID()
	if result.Status == module.StatusLogout {
		result.Message = "登录状态已失效，请重新登录"
		reply(w, http.StatusUnauthorized, result, nil)
		return
	}
}

func GetResult(w http.ResponseWriter, r *http.Request) *api.ReplyResult {
	result := &api.ReplyResult{Status: module.StatusLogout}
	userData, err := module.NewSession(w, r)
	if err == nil {
		result.UserID, result.Status = userData.UserID()
	}
	return result
}
