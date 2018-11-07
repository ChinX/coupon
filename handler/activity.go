package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/chinx/cobweb"

	"github.com/chinx/coupon/model"
)

func CreateActivity(w http.ResponseWriter, r *http.Request) {

}

func ListActivities(w http.ResponseWriter, r *http.Request) {
	byteData, err := pagedQuery(r, &model.Activity{})
	if err != nil {
		reply(w, http.StatusInternalServerError, "服务器开了会儿小差，请稍后尝试")
		return
	}
	reply(w, http.StatusOK, byteData)
}

func ModifyActivity(w http.ResponseWriter, r *http.Request) {

}

func GetActivity(w http.ResponseWriter, r *http.Request) {
	params := r.Context().Value(context.Background())
	if params == nil {
		reply(w, http.StatusBadRequest, "请求参数错误")
		return
	}

	id, err := strconv.Atoi(params.(cobweb.Params).Get("activity_id"))
	if err != nil || id == 0 {
		reply(w, http.StatusBadRequest, "请求参数错误")
		return
	}

	activity := &model.Activity{ID: int64(id)}
	if ok := model.Get(activity); !ok {
		reply(w, http.StatusBadRequest, "指定的活动不存在")
		return
	}

	reply(w, http.StatusOK, activity)
}

func DeleteActivity(w http.ResponseWriter, r *http.Request) {

}
