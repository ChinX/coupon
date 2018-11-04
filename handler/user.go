package handler

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/chinx/coupon/model"
	"github.com/chinx/coupon/service"
	"github.com/chinx/mohist/iorw"
	"github.com/go-session/session"
)

type authCode struct {
	Code string `json:"code"`
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	auth := &authCode{}
	err := iorw.ReadJSON(r.Body, auth)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("请求参数错误"))
		return
	}

	wxSession, err := service.NewAuth(auth.Code).AuthSession()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("从微信获取登录信息失败"))
		return
	}

	store, err := session.Start(context.Background(), w, r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("服务器开了会儿小差，请稍后尝试"))
		return
	}

	store.Set("openid", wxSession.OpenID)
	store.Set("session_key", wxSession.SessionKey)
	err = store.Save()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("正拿肇事程序员祭天，稍后片刻"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

type binding struct {
	Signature     string `json:"signature"`
	RawData       string `json:"rawData"`
	EncryptedData string `json:"encryptedData"`
	IV            string `json:"iv"`
}

func UserBinding(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(context.Background(), w, r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("连接超时，请重新登录"))
		return
	}

	b := &binding{}
	err = iorw.ReadJSON(r.Body, b)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("请求参数错误"))
		return
	}

	key, ok := store.Get("session_key")
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("请求参数错误"))
		return
	}

	sign := signature(b.RawData + key.(string))
	if b.Signature != sign {
		log.Println(sign)
		log.Println(b.Signature)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("请求参数错误"))
		return
	}

	user := &model.User{}
	err = json.Unmarshal([]byte(b.RawData), user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("请求参数错误"))
		return
	}
	openid, _ := store.Get("openid")
	condition := &model.User{ID: openid.(string)}
	user.ID = condition.ID

	has := model.Get(condition)
	if !has {
		ok = model.Insert(user)
		if !ok {
			log.Println("insert")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("正拿肇事程序员祭天，稍后片刻"))
			return
		}
	} else {
		if condition.AvatarURL != user.AvatarURL ||
			condition.City != user.City ||
			condition.Province != user.Province ||
			condition.Country != user.Country ||
			condition.Gender != user.Gender ||
			condition.Language != user.Language ||
			condition.Nickname != user.Nickname {
			ok = model.Update(user)
			if !ok {
				log.Println("Update")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("正拿肇事程序员祭天，稍后片刻"))
				return
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func UserLogout(w http.ResponseWriter, r *http.Request) {
	err := session.Destroy(context.Background(), w, r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("未登录"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func signature(key string) string {
	h := sha1.New()
	io.WriteString(h, key)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func CheckSession(w http.ResponseWriter, r *http.Request) {
	_, err := session.Start(context.Background(), w, r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("连接超时，请重新登录"))
		return
	}
	session.Refresh(context.Background(), w, r)
}
