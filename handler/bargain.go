package handler

import (
	"context"
	"github.com/chinx/coupon/model"
	"github.com/chinx/mohist/iorw"
	"github.com/go-session/session"
	"log"
	"net/http"
)

func CreateBargain(w http.ResponseWriter, r *http.Request) {
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
	store.SessionID()
	userID, ok := store.Get("openid")
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("请求参数错误"))
		return
	}

	activity := &model.Activity{ID: t.ActivityID}
	if ok := model.Get(activity); !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("activity id is wrong"))
		return
	}

	mTask := &model.Task{
		UserID: userID.(string),
		ActivityID: t.ActivityID,
	}

	if ok := model.Get(mTask); ok {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte("不能重复领取任务"))
		return
	}

	mTask.Message= t.Message
	mTask.Price = activity.Price

	if ok := model.Insert(mTask); !ok{
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("activity id is wrong"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func ListBargains(w http.ResponseWriter, r *http.Request) {

}
