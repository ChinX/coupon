package handler

import (
	"github.com/chinx/coupon/model"
	"net/http"
)

type bargainMsg struct {
	Message string `json:"message"`
	TaskID  int64  `json:"task_id"`
}

func CreateBargain(w http.ResponseWriter, r *http.Request) {
	userID, err := readUserID(w, r)
	if err != nil {
		reply(w, http.StatusUnauthorized, "登录状态已失效，请重新登录")
		return
	}

	msg := &bargainMsg{}
	err = readBody(r.Body, msg)
	if err != nil {
		reply(w, http.StatusInternalServerError, "服务器开了会儿小差，请稍后尝试")
		return
	}

	task := &model.Task{ID: msg.TaskID}
	if ok := model.Get(task); !ok {
		reply(w, http.StatusBadRequest, "指定的任务不存在")
		return
	}

	bargain := &model.Bargain{
		UserID: userID,
		TaskID: msg.TaskID,
	}

	if ok := model.Get(bargain); ok {
		reply(w, http.StatusConflict, "不能重复砍刀")
		return
	}

	count := discount(task.Price, 5)
	task.Progress += count

	if task.Progress >= float64(task.Price){
		task.Progress = float64(task.Price)
	}
	bargain.Message = msg.Message
	bargain.Discount = count

	session := model.NewSession()
	defer session.Close()

	err = session.Begin()
	_, err = session.Insert(bargain)
	if err != nil {
		session.Rollback()
		reply(w, http.StatusInternalServerError, err)
		return
	}

	_, err = session.Where("id = ?", msg.TaskID).Update(task)
	if err != nil {
		session.Rollback()
		reply(w, http.StatusInternalServerError, err)
		return
	}

	if task.Progress == float64(task.Price){
		coupon := &model.Coupon{
			UserID: task.UserID,
			ActivityID: task.ActivityID,
		}
		_, err = session.Insert(coupon)
		if err != nil {
			session.Rollback()
			reply(w, http.StatusInternalServerError, err)
			return
		}

		//_, err = session.Where("id = ?", task.ActivityID).Update(Activity)
		//if err != nil {
		//	session.Rollback()
		//	return
		//}

	}

	err = session.Commit()
	if err != nil {
		reply(w, http.StatusInternalServerError, err)
		return
	}

	reply(w, http.StatusCreated, bargain)
}

func ListBargains(w http.ResponseWriter, r *http.Request) {
	byteData, err := pagedQuery(r, &model.Bargain{})
	if err != nil {
		reply(w, http.StatusInternalServerError, "服务器开了会儿小差，请稍后尝试")
		return
	}
	reply(w, http.StatusOK, byteData)
}

func discount(total int64, share int64) float64 {
	return 1
}
