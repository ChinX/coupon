package service

import (
	"fmt"
	"github.com/chinx/mohist/iorw"
	"net/http"
	"net/url"
)

var (
	appID     = "wx3f73a5186ad1702a"
	appSecret = "d11cf3f8b7ea0ee37100be01431e289a"
	authURL   = "https://api.weixin.qq.com/sns/jscode2session"
)

type WXAuth struct {
	authArgs url.Values
}

type WXSession struct {
	SessionKey string `json:"session_key"`
	OpenID     string `json:"openid"`
}

func NewAuth(code string) *WXAuth {
	args := url.Values{}
	args.Set("appid", appID)
	args.Set("secret", appSecret)
	args.Set("js_code", code)
	args.Set("grant_type", "authorization_code")
	return &WXAuth{
		authArgs: args,
	}
}

func (wx *WXAuth) AuthSession() (*WXSession, error) {
	resp, err := http.Get(authURL + "?" + wx.authArgs.Encode())
	if err != nil || resp.StatusCode != 200 {
		return nil, fmt.Errorf("get auth session from wexin faild: %s", err)
	}

	session := &WXSession{}
	err = iorw.ReadJSON(resp.Body, session)
	if err != nil {
		return nil, fmt.Errorf("parse jsom body error: %s", err)
	}

	return session, nil
}
