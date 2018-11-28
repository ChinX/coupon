package handler

import (
	"log"
	"net/http"

	"github.com/chinx/coupon/api"
	"github.com/chinx/coupon/module"
)

func UserLogin(w http.ResponseWriter, r *http.Request) {
	operation := "UserLogin"
	log.Println(operation, r.Header.Get("Cookie"))
	auth := &api.UserLogin{}
	result := &api.ReplyResult{Status: module.StatusLogout}
	err := readBody(r.Body, auth)
	if err != nil {
		result.Message = "请求参数错误"
		log.Println(operation, err)
		reply(w, http.StatusBadRequest, result, err)
		return
	}

	log.Println(operation, *auth)
	wxData, err := module.NewAuth(auth.Code).AuthSession()
	if err != nil {
		log.Println(operation, err)
		result.Message = "从微信获取登录信息失败"
		reply(w, http.StatusBadRequest, result, err)
		return
	}

	log.Println(operation, *wxData)
	userData, err := module.NewSession(w, r)
	if err != nil {
		log.Println(operation, err)
		result.Message = "创建登录信息失败"
		reply(w, http.StatusUnauthorized, result, err)
		return
	}

	if userDI, _ := userData.UserID(); userDI != "" && userDI != wxData.OpenID{
		log.Printf("%s != %s, delete and refresh cookie\n", userDI, wxData.OpenID)
		userData.Refresh(w, r)
	}

	userData.ShowALL()
	result.UserID = wxData.OpenID
	result.Status = userData.SetUserSession(wxData)
	if result.Status == module.StatusLogout {
		log.Println(operation, err)
		result.Message = "登录失败，请稍后尝试"
		reply(w, http.StatusUnauthorized, result, nil)
		return
	}

	if result.Status == module.StatusLogin{
		user, ok := module.GetUserInfo(result.UserID)
		log.Println(operation, user, ok)
		if ok {
			result.Data = user
		}
	}

	reply(w, http.StatusCreated, result, nil)
}

func UserBinding(w http.ResponseWriter, r *http.Request) {
	operation := "UserBinding"
	log.Println(operation, r.Header.Get("Cookie"))
	result := &api.ReplyResult{Status: module.StatusLogout}
	userData, err := module.NewSession(w, r)
	if err != nil {
		log.Println(operation, err)
		result.Message = "获取登录信息失败"
		reply(w, http.StatusUnauthorized, result, err)
		return
	}

	userData.ShowALL()

	bind := &module.Binding{}
	err = readBody(r.Body, bind)
	if err != nil {
		log.Println(operation, err)
		result.Message = "请求参数错误"
		reply(w, http.StatusBadRequest, result, err)
		return
	}

	log.Println(operation, *bind)
	status, err := userData.ValidateSignature(bind.Signature, bind.RawData)
	result.Status = status
	if err != nil {
		log.Println(operation, err)
		result.Message = err.Error()
		reply(w, http.StatusBadRequest, result, err)
		return
	}

	wxUser, err := userData.Binding(bind)
	if err != nil {
		log.Println(operation, err)
		result.Message = err.Error()
		reply(w, http.StatusBadRequest, result, err)
		return
	}

	result.UserID, result.Data = wxUser.ID, wxUser
	result.Status = module.StatusLogin
	reply(w, http.StatusCreated, result, nil)
}

func checkUser(w http.ResponseWriter, r *http.Request) *api.ReplyResult {
	return checkLogin(w, r, module.PermissionUser)
}

func checkAdmin(w http.ResponseWriter, r *http.Request) *api.ReplyResult {
	return checkLogin(w, r, module.PermissionAdmin)
}
