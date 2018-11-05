package handler

import (
	"context"
	"encoding/json"
	"github.com/chinx/cobweb"
	"github.com/chinx/coupon/model"
	"log"
	"net/http"
	"strconv"

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

func ListTasks(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	from, err := strconv.Atoi(r.Form.Get("from"))
	if err != nil || from < 0 {
		from = 0
	}

	count, err := strconv.Atoi(r.Form.Get("count"))
	if err != nil || count < 30 {
		count = 30
	}

	mTask := &model.Task{}
	n, list := mTask.List(from, count)
	result := &listResult{
		Total: n,
		List:  list,
	}

	byteData, err := json.Marshal(result)
	if err != nil {
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("正拿肇事程序员祭天，稍后片刻"))
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write(byteData)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	params:=r.Context().Value(context.Background())
	if params == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("activity id is not found"))
		return
	}

	id, err := strconv.Atoi(params.(cobweb.Params).Get("task_id"))
	if err != nil || id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("task id is wrong"))
		return
	}

	mTask := &model.Task{ID: int64(id)}
	if ok := model.Get(mTask); !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("activity id is wrong"))
		return
	}

	byteData, err := json.Marshal(mTask)
	if err != nil {
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("正拿肇事程序员祭天，稍后片刻"))
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write(byteData)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {

}
