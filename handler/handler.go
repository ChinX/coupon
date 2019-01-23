package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/chinx/coupon/module"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://www.windup.cn", http.StatusMovedPermanently)
}

func VerificationHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("2efb1213bd8a6bd6049be787f2480189"))
}

func ListActivities(w http.ResponseWriter, r *http.Request) {
	operation := "[ListActivities]"
	log.Println(operation)
	result := GetResult(w, r)
	params := pageParams(r)
	result.Data = module.ListActivities(result.UserID, params.From, params.Count)
	reply(w, http.StatusOK, result, nil)
}

func UserTask(w http.ResponseWriter, r *http.Request) {
	operation := "[UserTask]"
	log.Println(operation)
	result := GetResult(w, r)
	userID, err := strconv.Atoi(urlParam(r, "user_id"))
	if err != nil || userID == 0 {
		log.Println(operation, err)
		result.Message = "请求参数错误"
		reply(w, http.StatusBadRequest, result, nil)
		return
	}
	activityID, err := strconv.Atoi(urlParam(r, "activity_id"))
	if err != nil || activityID == 0 {
		log.Println(operation, err)
		result.Message = "请求参数错误"
		reply(w, http.StatusBadRequest, result, nil)
		return
	}

	data, err := module.UserTask(result.UserID, int64(userID), int64(activityID))
	if err != nil {
		log.Println(operation, err)
		result.Message = err.Error()
		reply(w, http.StatusInternalServerError, result, err)
		return
	}
	result.Data = data
	reply(w, http.StatusOK, result, err)
}

func ListBargain(w http.ResponseWriter, r *http.Request) {
	operation := "[ListBargain]"
	log.Println(operation)
	result := GetResult(w, r)

	taskID, err := strconv.Atoi(urlParam(r, "task_id"))
	if err != nil || taskID == 0 {
		log.Println(operation, err)
		result.Message = "请求参数错误"
		reply(w, http.StatusBadRequest, result, nil)
		return
	}

	params := pageParams(r)
	result.Data = module.ListBargains(int64(taskID), params.From, params.Count)
	reply(w, http.StatusOK, result, nil)
}

func CreateBargain(w http.ResponseWriter, r *http.Request) {
	operation := "[CreateBargain]"
	result := GetResult(w, r)
	taskID, err := strconv.Atoi(urlParam(r, "task_id"))
	if err != nil || taskID == 0 {
		log.Println(operation, err)
		result.Message = "请求参数错误"
		reply(w, http.StatusBadRequest, result, nil)
		return
	}

	taskBargain, err := module.CreateBargain(result.UserID, int64(taskID))
	if err != nil {
		log.Println(operation, err)
		result.Message = err.Error()
		reply(w, http.StatusInternalServerError, result, err)
		return
	}
	result.Data = taskBargain
	reply(w, http.StatusCreated, result, nil)
}

func CreateCash(w http.ResponseWriter, r *http.Request) {
	operation := "[CreateCash]"
	result := GetResult(w, r)
	taskID, err := strconv.Atoi(urlParam(r, "task_id"))
	if err != nil || taskID == 0 {
		log.Println(operation, err)
		result.Message = "请求参数错误"
		reply(w, http.StatusBadRequest, result, nil)
		return
	}

	task, err := module.CreateCash(result.UserID, int64(taskID))
	if err != nil {
		log.Println(operation, err)
		result.Message = err.Error()
		reply(w, http.StatusInternalServerError, result, err)
		return
	}

	result.Data = task
	reply(w, http.StatusCreated, task, nil)
}

var StaticDir = "./static"

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = urlParam(r, "filename")
	http.FileServer(http.Dir(StaticDir)).ServeHTTP(w, r)
}
