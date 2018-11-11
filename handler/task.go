package handler

import (
	"net/http"
	"strconv"

	"github.com/chinx/coupon/api"
	"github.com/chinx/coupon/module"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {
	result := checkUser(w, r)
	if result.Status != module.StatusLogin {
		return
	}

	activityID, err := strconv.Atoi(urlParam(r, "activity_id"))
	if err != nil || activityID == 0 {
		result.Message = "请求参数错误"
		reply(w, http.StatusBadRequest, result, nil)
		return
	}

	msg := &api.Message{}
	err = readBody(r.Body, msg)
	if err != nil {
		result.Message = "请求参数错误"
		reply(w, http.StatusBadRequest, result, err)
		return
	}

	task, err := module.CreateTask(result.UserID, int64(activityID), msg.Message)
	if err != nil {
		result.Message = err.Error()
		reply(w, http.StatusInternalServerError, result, err)
		return
	}

	taskBargain, err := module.CreateBargain(result.UserID, task.ID, msg.Message)
	if err != nil {
		result.Message = "砍刀失败"
		reply(w, http.StatusBadRequest, result, err)
		return
	}

	reply(w, http.StatusCreated, taskBargain, nil)
}

func ActiveTask(w http.ResponseWriter, r *http.Request) {
	result := checkUser(w, r)
	if result.Status != module.StatusLogin {
		return
	}

	activityID, err := strconv.Atoi(urlParam(r, "activity_id"))
	if err != nil || activityID == 0 {
		result.Message = "请求参数错误"
		reply(w, http.StatusBadRequest, result, nil)
		return
	}

	task, err := module.GetTaskByUserActivity(result.UserID, int64(activityID))
	if err != nil {
		result.Message = err.Error()
		reply(w, http.StatusBadRequest, result, err)
		return
	}

	reply(w, http.StatusOK, task, nil)

}

func ListTasks(w http.ResponseWriter, r *http.Request) {
	result := checkUser(w, r)
	if result.Status != module.StatusLogin {
		return
	}

	params := pageParams(r)
	byteData, err := pagedResult(module.ListTasks(result.UserID, params.From, params.Count))
	if err != nil {
		result.Message = "获取任务列表失败"
		reply(w, http.StatusInternalServerError, result, err)
		return
	}
	reply(w, http.StatusOK, byteData, nil)
}

func GetTask(w http.ResponseWriter, r *http.Request) {
	result := checkUser(w, r)
	if result.Status != module.StatusLogin {
		return
	}

	taskID, err := strconv.Atoi(urlParam(r, "task_id"))
	if err != nil || taskID == 0 {
		result.Message = "请求参数错误"
		reply(w, http.StatusBadRequest, result, nil)
		return
	}

	task, err := module.GetTask(int64(taskID))
	if err != nil {
		result.Message = err.Error()
		reply(w, http.StatusBadRequest, result, err)
		return
	}

	reply(w, http.StatusOK, task, nil)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	result := checkUser(w, r)
	if result.Status != module.StatusLogin {
		return
	}

	taskID, err := strconv.Atoi(urlParam(r, "task_id"))
	if err != nil || taskID == 0 {
		result.Message = "请求参数错误"
		reply(w, http.StatusBadRequest, result, nil)
		return
	}

	task, err := module.DeleteTask(int64(taskID))
	if err != nil {
		result.Message = err.Error()
		reply(w, http.StatusBadRequest, result, err)
		return
	}

	reply(w, http.StatusCreated, task, nil)
}
