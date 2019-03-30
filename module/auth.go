package module

import (
	"errors"
	"fmt"
	"github.com/chinx/mohist/iorw"
	"net/http"
	"net/url"

	"github.com/chinx/coupon/dao/mysql"
	"github.com/chinx/coupon/model"
)

var (
	AppID     = "wx3f73a5186ad1702a"
	AppSecret = "3a30eacdb77a59ad4b281c1daee38e3d"
	authURL   = "https://api.weixin.qq.com/sns/jscode2session"
)

type WXAuth struct {
	authArgs url.Values
}

type WXSession struct {
	ID         int64  `json:"-"`
	SessionKey string `json:"session_key"`
	OpenID     string `json:"openid"`
}

type Binding struct {
	Signature     string `json:"signature"`
	RawData       string `json:"rawData"`
	EncryptedData string `json:"encryptedData"`
	IV            string `json:"iv"`
}

func NewAuth(code string) *WXAuth {
	args := url.Values{}
	args.Set("appid", AppID)
	args.Set("secret", AppSecret)
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
	if err != nil || session.OpenID == "" || session.SessionKey == ""{
		return nil, fmt.Errorf("parse jsom body error: %s", err)
	}

	user := &model.User{}

	err = mysql.Get(user, "openid=?", session.OpenID)
	if err == mysql.NoRecords {
		user.Openid = session.OpenID
		err = mysql.Insert(user)
	}
	if err != nil {
		return nil, errors.New("绑定用户信息失败")
	}
	session.ID = user.ID
	return session, nil
}
