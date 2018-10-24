package handler

import (
	"fmt"
	"github.com/naoina/denco"
	"io/ioutil"
	"net/http"
	"net/url"
)

var (
	appID string = "wx3f73a5186ad1702a"
	appSecret = "d11cf3f8b7ea0ee37100be01431e289a"
	authURL = "https://api.weixin.qq.com/sns/jscode2session"
)

func NewWxAuth(code string) url.Values {
	auth := url.Values{}
	auth.Set("appid", appID)
	auth.Set("secret", appSecret)
	auth.Set("js_code", code)
	auth.Set("grant_type", "authorization_code")
	return auth
}

func UserLogin(w http.ResponseWriter, r *http.Request, params denco.Params) {
	r.ParseForm()
	callURL := authURL+ "?" + NewWxAuth(r.Form.Get("code")).Encode()
	fmt.Println(callURL)
	resp, err := http.Get(callURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	byteArr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	fmt.Println(string(byteArr))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
