package handler

import (
	"context"
	"encoding/json"
	"github.com/chinx/cobweb"
	"github.com/chinx/coupon/api"
	"github.com/chinx/coupon/module"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type pagination struct {
	Total int64       `json:"total"`
	List  interface{} `json:"list"`
}

func pageParams(r *http.Request) *api.PageParams {
	r.ParseForm()
	params := &api.PageParams{}
	params.From, _ = strconv.Atoi(r.Form.Get("from"))
	if params.From < 1 {
		params.From = 1
	}

	params.Count, _ = strconv.Atoi(r.Form.Get("count"))
	if params.Count < 0 || params.Count > 100 {
		params.Count = 30
	}

	log.Println("pageParams", *params)
	return params
}

func urlParam(r *http.Request, key string) string {
	params := r.Context().Value(context.Background())
	if params == nil {
		return ""
	}
	log.Println("urlParam", params)
	return params.(cobweb.Params).Get(key)
}

func pagedResult(total int64, list interface{}) ([]byte, error) {
	return json.Marshal(&pagination{
		Total: total,
		List:  list,
	})
}

func readBody(body io.Reader, v interface{}) error {
	byt, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	log.Println(string(byt))
	err = json.Unmarshal(byt, v)
	if err != nil {
		return err
	}
	return nil
}

func reply(w http.ResponseWriter, status int, data interface{}, err error) {
	var result []byte
	switch t := data.(type) {
	case []byte:
		result = t
	case string:
		result = []byte(t)
	case error:
		result = []byte(t.Error())
	default:
		byteData, err := json.Marshal(data)
		if err != nil {
			status = http.StatusInternalServerError
			result = []byte(err.Error())
		} else {
			result = byteData
		}
	}

	log.Println(string(result))
	if status >= http.StatusBadRequest {
		if err != nil {
			log.Println(err)
		} else {


			log.Println(status, string(result))
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(result)
}

func checkLogin(w http.ResponseWriter, r *http.Request, permission int) *api.ReplyResult {
	result := &api.ReplyResult{Status: module.StatusLogout}
	userData, err := module.NewSession(w, r)
	if err != nil {
		log.Println("checkLogin", err)
		result.Message = "获取登录信息失败"
		reply(w, http.StatusUnauthorized, result, err)
		return result
	}

	log.Println("checkLogin")
	userData.ShowALL()
	if !userData.IsPermission(permission) {
		result.Message = "无权限"
		reply(w, http.StatusForbidden, result, err)
		return result
	}

	result.UserID, result.Status = userData.UserID()
	if result.Status == module.StatusLogout {
		result.Message = "登录状态已失效，请重新登录"
		reply(w, http.StatusUnauthorized, result, nil)
		return result
	}
	return result
}
