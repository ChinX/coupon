package handler

import (
	"net/http"

	"github.com/chinx/coupon/api"
	"github.com/chinx/coupon/module"
)

func UserLogin(w http.ResponseWriter, r *http.Request) {
	auth := &api.UserLogin{}
	result := &api.ReplyResult{Status: module.StatusLogout}
	err := readBody(r.Body, auth)
	if err != nil {
		result.Message = "请求参数错误"
		reply(w, http.StatusBadRequest, result, err)
		return
	}

	wxData, err := module.NewAuth(auth.Code).AuthSession()
	if err != nil {
		result.Message = "从微信获取登录信息失败"
		reply(w, http.StatusBadRequest, result, err)
		return
	}

	userData, err := module.NewSession(w, r)
	if err != nil {
		result.Message = "创建登录信息失败"
		reply(w, http.StatusUnauthorized, result, err)
		return
	}

	result.Status = userData.SetUserSession(wxData)
	if result.Status == module.StatusLogout {
		result.Message = "登录失败，请稍后尝试"
		reply(w, http.StatusUnauthorized, result, nil)
		return
	}

	if result.Status == module.StatusLogin{
		user, ok := module.GetUserInfo(result.UserID)
		if ok {
			result.Data = user
		}
	}

	reply(w, http.StatusCreated, result, nil)
}

func UserBinding(w http.ResponseWriter, r *http.Request) {
	result := &api.ReplyResult{Status: module.StatusLogout}
	userData, err := module.NewSession(w, r)
	if err != nil {
		result.Message = "获取登录信息失败"
		reply(w, http.StatusUnauthorized, result, err)
		return
	}

	bind := &module.Binding{}
	err = readBody(r.Body, bind)
	if err != nil {
		result.Message = "请求参数错误"
		reply(w, http.StatusBadRequest, result, err)
		return
	}

	status, err := userData.ValidateSignature(bind.Signature, bind.RawData)
	result.Status = status
	if err != nil {
		result.Message = err.Error()
		reply(w, http.StatusBadRequest, result, err)
		return
	}

	err = userData.Binding(bind)
	if err != nil {
		result.Message = err.Error()
		reply(w, http.StatusBadRequest, result, err)
		return
	}

	result.Data, _ = module.GetUserInfo(result.UserID)
	result.Status = module.StatusLogin
	reply(w, http.StatusCreated, result, nil)
}

func checkUser(w http.ResponseWriter, r *http.Request) *api.ReplyResult {
	return checkLogin(w, r, module.PermissionUser)
}

func checkAdmin(w http.ResponseWriter, r *http.Request) *api.ReplyResult {
	return checkLogin(w, r, module.PermissionAdmin)
}
