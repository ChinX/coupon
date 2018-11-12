package module

import (
	"errors"
	"github.com/chinx/coupon/model"
	"time"
)

func CreateCoupon(userID string, activityID int64)(*model.Coupon, error) {
	task, err := GetTaskByUserActivity(userID, activityID)
	if err != nil {
		return nil, err
	}

	if task.Status != model.TaskStatusDone{
		return nil, errors.New("未完成相关任务")
	}

	coupon := &model.Coupon{
		UserID: userID,
		ActivityID: activityID,
	}

	if ok := model.Get(coupon); ok {
		return nil, errors.New("不能重复领取门票")
	}

	activity, err := GetActiveDetail(activityID)
	if err != nil {
		return nil, err
	}

	if activity.ExpireDate.Before(time.Now()){
		return nil, errors.New("你来晚了，活动已失效")
	}

	if activity.Total <= activity.Completed {
		return nil, errors.New("你来晚了，任务门票已被领取完")
	}


	if activity.DailyTotal <= activity.DailyCompleted{
		return nil, errors.New("你来晚了，今天任务门票已被领取完")
	}

	if ok := model.Insert(coupon); !ok{
		return nil, errors.New("领取门票失败")
	}

	return coupon, nil
}

func ListCoupons(userID string, from, count int) (int64, interface{}) {
	coupon := &model.Coupon{UserID: userID}
	return coupon.List(from, count)
}


func GetCoupon(id int64) (*model.Coupon, error) {
	coupon := &model.Coupon{ID: id}
	if ok := model.Get(coupon); !ok {
		return nil, errors.New( "指定的门票不存在")
	}
	return coupon, nil
}

func GetCouponByActivity(userID string, activityID int64) (*model.Coupon, error) {
	coupon := &model.Coupon{UserID: userID, ActivityID: activityID}
	if ok := model.Get(coupon); !ok {
		return nil, errors.New( "指定的门票不存在")
	}
	return coupon, nil
}

func DeleteCoupon(id int64) (*model.Coupon, error) {
	coupon := &model.Coupon{ID: id}
	if ok := model.Delete(coupon); !ok {
		return nil, errors.New( "兑换实体门票失败")
	}
	return coupon, nil
}