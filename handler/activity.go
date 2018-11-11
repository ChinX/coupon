package handler

import (
	"net/http"
	"strconv"

	"github.com/chinx/coupon/module"
)

func ListActivities(w http.ResponseWriter, r *http.Request) {
	result := checkUser(w, r)
	if result.Status != module.StatusLogin {
		return
	}

	params := pageParams(r)
	byteData, err := pagedResult(module.ListActivities(params.From, params.Count))
	if err != nil {
		result.Message = "获取活动列表失败"
		reply(w, http.StatusInternalServerError, result, err)
		return
	}
	reply(w, http.StatusOK, byteData, nil)
}

func GetActivity(w http.ResponseWriter, r *http.Request) {
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

	activity, err := module.GetActivity(int64(activityID))
	if err != nil {
		result.Message = "指定的活动不存在"
		reply(w, http.StatusBadRequest, result, nil)
		return
	}

	reply(w, http.StatusOK, activity, nil)
}

func CreateActivity(w http.ResponseWriter, r *http.Request) {
	result := checkAdmin(w, r)
	if result.Status != module.StatusLogin {
		return
	}
}

func ModifyActivity(w http.ResponseWriter, r *http.Request) {
	result := checkAdmin(w, r)
	if result.Status != module.StatusLogin {
		return
	}

	id, err := strconv.Atoi(urlParam(r, "activity_id"))
	if err != nil || id == 0 {
		result.Message = "请求参数错误"
		reply(w, http.StatusBadRequest, result, nil)
		return
	}
	reply(w, http.StatusOK, "ok", nil)
}

func DeleteActivity(w http.ResponseWriter, r *http.Request) {
	result := checkAdmin(w, r)
	if result.Status != module.StatusLogin {
		return
	}

	activityID, err := strconv.Atoi(urlParam(r, "activity_id"))
	if err != nil || activityID == 0 {
		result.Message = "请求参数错误"
		reply(w, http.StatusBadRequest, result, nil)
		return
	}

	activity, err := module.DeleteActivity(int64(activityID))
	if err != nil {
		result.Message = "指定的活动不存在"
		reply(w, http.StatusBadRequest, result, nil)
		return
	}
	reply(w, http.StatusOK, activity, nil)
}
