package handler

import (
	"context"
	"github.com/chinx/cobweb"
	"net/http"
	"strconv"

	"github.com/chinx/coupon/model"
)

func ListCoupons(w http.ResponseWriter, r *http.Request) {
	byteData, err := pagedQuery(r, &model.Coupon{})
	if err != nil {
		reply(w, http.StatusInternalServerError, "服务器开了会儿小差，请稍后尝试")
		return
	}
	reply(w, http.StatusOK, byteData)
}

func ModifyCoupon(w http.ResponseWriter, r *http.Request) {

}

func GetCoupon(w http.ResponseWriter, r *http.Request) {
	params := r.Context().Value(context.Background())
	if params == nil {
		reply(w, http.StatusBadRequest, "请求参数错误")
		return
	}

	id, err := strconv.Atoi(params.(cobweb.Params).Get("coupon_id"))
	if err != nil || id == 0 {
		reply(w, http.StatusBadRequest, "请求参数错误")
		return
	}

	coupon := &model.Coupon{ID: int64(id)}
	if ok := model.Get(coupon); !ok {
		reply(w, http.StatusBadRequest, "指定的优惠券不存在")
		return
	}

	reply(w, http.StatusOK, coupon)
}
