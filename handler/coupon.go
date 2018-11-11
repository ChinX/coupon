package handler

import (
	"net/http"
	"strconv"

	"github.com/chinx/coupon/module"
)

func CreateCoupon(w http.ResponseWriter, r *http.Request) {
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

	coupon, err := module.CreateCoupon(result.UserID, int64(activityID))
	if err != nil {
		result.Message = err.Error()
		reply(w, http.StatusInternalServerError, result, err)
		return
	}

	reply(w, http.StatusCreated, coupon, nil)
}

func ActiveCoupon(w http.ResponseWriter, r *http.Request) {
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

	coupon, err := module.GetCouponByActivity(result.UserID, int64(activityID))
	if err != nil {
		result.Message = err.Error()
		reply(w, http.StatusInternalServerError, result, err)
		return
	}

	reply(w, http.StatusOK, coupon, nil)

}

func ListCoupons(w http.ResponseWriter, r *http.Request) {
	result := checkUser(w, r)
	if result.Status != module.StatusLogin {
		return
	}

	params := pageParams(r)
	byteData, err := pagedResult(module.ListCoupons(result.UserID, params.From, params.Count))
	if err != nil {
		result.Message = "获取门票列表失败"
		reply(w, http.StatusInternalServerError, result, err)
		return
	}
	reply(w, http.StatusOK, byteData, nil)
}

func GetCoupon(w http.ResponseWriter, r *http.Request) {
	result := checkUser(w, r)
	if result.Status != module.StatusLogin {
		return
	}

	couponID, err := strconv.Atoi(urlParam(r, "coupon_id"))
	if err != nil || couponID == 0 {
		result.Message = "请求参数错误"
		reply(w, http.StatusBadRequest, result, nil)
		return
	}

	coupon, err := module.GetCoupon(int64(couponID))
	if err != nil {
		result.Message = err.Error()
		reply(w, http.StatusBadRequest, result, nil)
		return
	}

	reply(w, http.StatusOK, coupon, nil)
}

func DeleteCoupon(w http.ResponseWriter, r *http.Request) {
	result := checkUser(w, r)
	if result.Status != module.StatusLogin {
		return
	}

	couponID, err := strconv.Atoi(urlParam(r, "coupon_id"))
	if err != nil || couponID == 0 {
		result.Message = "请求参数错误"
		reply(w, http.StatusBadRequest, result, nil)
		return
	}

	coupon, err := module.DeleteCoupon(int64(couponID))
	if err != nil {
		result.Message = err.Error()
		reply(w, http.StatusBadRequest, result, nil)
		return
	}

	reply(w, http.StatusCreated, coupon, nil)
}
