package handler

import (
	"context"
	"github.com/chinx/coupon/model"
	"log"
	"net/http"

	"github.com/chinx/mohist/iorw"
	"github.com/go-session/session"
)

type task struct {
	Message    string `json:"message"`
	ActivityID int64  `json:"activity_id"`
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	t := &task{}
	err := iorw.ReadJSON(r.Body, t)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("请求参数错误"))
		return
	}
	store, err := session.Start(context.Background(), w, r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("服务器开了会儿小差，请稍后尝试"))
		return
	}
	userID, ok := store.Get("opendi")
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("请求参数错误"))
		return
	}

	mTask := &model.Task{
		UserID: userID.(string),
		Message: t.Message,
		ActivityID: t.ActivityID,
	}

	t.UserID = userID.(string)

	store.Set("openid", wxSession.OpenID)
}

func ListTasks(w http.ResponseWriter, r *http.Request) {

}

func GetTask(w http.ResponseWriter, r *http.Request) {

}

func DeleteTask(w http.ResponseWriter, r *http.Request) {

}
