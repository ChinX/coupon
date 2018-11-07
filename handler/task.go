package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/chinx/cobweb"
	"github.com/chinx/coupon/model"
)

type taskMsg struct {
	Message    string `json:"message"`
	ActivityID int64  `json:"activity_id"`
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	userID, err := readUserID(w, r)
	if err != nil {
		reply(w, http.StatusUnauthorized, "登录状态已失效，请重新登录")
		return
	}

	msg := &taskMsg{}
	err = readBody(r.Body, msg)
	if err != nil {
		reply(w, http.StatusInternalServerError, "服务器开了会儿小差，请稍后尝试")
		return
	}

	activity := &model.Activity{ID: msg.ActivityID}
	if ok := model.Get(activity); !ok {
		reply(w, http.StatusBadRequest, "指定的活动不存在")
		return
	}

	task := &model.Task{
		UserID:     userID,
		ActivityID: msg.ActivityID,
	}

	if ok := model.Get(task); ok {
		reply(w, http.StatusConflict, "不能重复领取任务")
		return
	}

	count := discount(task.Price, 5)
	task.Message = msg.Message
	task.Price = activity.Price
	task.Progress += count

	session := model.NewSession()
	defer session.Close()

	err = session.Begin()
	_, err = session.Insert(task)
	if err != nil {
		session.Rollback()
		reply(w, http.StatusInternalServerError, err)
		return
	}

	bargain := &model.Bargain{
		UserID:   userID,
		TaskID:   task.ID,
		Message:  "要对自己感情真，一刀砍到绝对深",
		Discount: count,
	}
	_, err = session.Insert(bargain)
	if err != nil {
		session.Rollback()
		reply(w, http.StatusInternalServerError, err)
		return
	}

	err = session.Commit()
	if err != nil {
		reply(w, http.StatusInternalServerError, err)
		return
	}
	reply(w, http.StatusCreated, task)
}

func ListTasks(w http.ResponseWriter, r *http.Request) {
	byteData, err := pagedQuery(r, &model.Task{})
	if err != nil {
		reply(w, http.StatusInternalServerError, "服务器开了会儿小差，请稍后尝试")
		return
	}
	reply(w, http.StatusOK, byteData)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	params := r.Context().Value(context.Background())
	if params == nil {
		reply(w, http.StatusBadRequest, "请求参数错误")
		return
	}

	id, err := strconv.Atoi(params.(cobweb.Params).Get("task_id"))
	if err != nil || id == 0 {
		reply(w, http.StatusBadRequest, "请求参数错误")
		return
	}

	task := &model.Task{ID: int64(id)}
	if ok := model.Get(task); !ok {
		reply(w, http.StatusBadRequest, "指定的任务ID不存在")
		return
	}

	reply(w, http.StatusOK, task)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {

}
