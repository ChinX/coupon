package handler

import (
	"context"
	"encoding/json"
	"github.com/chinx/cobweb"
	"log"
	"net/http"
	"strconv"

	"github.com/chinx/coupon/model"
)

type listResult struct {
	Total int64       `json:"total"`
	List  interface{} `json:"list"`
}

func CreateActivity(w http.ResponseWriter, r *http.Request) {

}

func ListActivities(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	from, err := strconv.Atoi(r.Form.Get("from"))
	if err != nil || from < 0 {
		from = 0
	}

	count, err := strconv.Atoi(r.Form.Get("count"))
	if err != nil || count < 30 {
		count = 30
	}

	activity := &model.Activity{}
	n, list := activity.List(from, count)
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

func ModifyActivity(w http.ResponseWriter, r *http.Request) {

}

func GetActivity(w http.ResponseWriter, r *http.Request) {
	params:=r.Context().Value(context.Background())
	if params == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("activity id is not found"))
		return
	}

	id, err := strconv.Atoi(params.(cobweb.Params).Get("activity_id"))
	if err != nil || id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("activity id is wrong"))
		return
	}

	activity := &model.Activity{ID: int64(id)}
	if ok := model.Get(activity); !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("activity id is wrong"))
		return
	}

	byteData, err := json.Marshal(activity)
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

func DeleteActivity(w http.ResponseWriter, r *http.Request) {

}
