package handler

import (
	"net/http"
	"strconv"

	"github.com/chinx/coupon/api"
	"github.com/chinx/coupon/module"
)

func CreateBargain(w http.ResponseWriter, r *http.Request) {
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

	msg := &api.Message{}
	err = readBody(r.Body, msg)
	if err != nil {
		result.Message = "请求参数错误"
		reply(w, http.StatusBadRequest, result, err)
		return
	}

	taskBargain, err := module.CreateBargain(result.UserID, int64(taskID), msg.Message)
	if err != nil {
		result.Message = "砍刀失败"
		reply(w, http.StatusBadRequest, result, err)
		return
	}

	reply(w, http.StatusCreated, taskBargain, nil)
}

func ListBargains(w http.ResponseWriter, r *http.Request) {
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
	params := pageParams(r)
	byteData, err := pagedResult(module.ListBargains(int64(taskID), params.From, params.Count))
	if err != nil {
		result.Message = "获取砍刀列表失败"
		reply(w, http.StatusInternalServerError, result, err)
		return
	}
	reply(w, http.StatusOK, byteData, nil)
}
