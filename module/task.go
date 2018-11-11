package module

import (
	"errors"

	"github.com/chinx/coupon/model"
)

func CreateTask(userID string, activityID int64, msg string) (*model.Task, error) {
	detail := &model.ActiveDetail{ID: activityID}
	if ok := model.Get(detail); !ok {
		return nil, errors.New("指定任务不存在")
	}

	task := &model.Task{
		UserID:     userID,
		ActivityID: activityID,
	}

	if ok := model.Get(task); ok {
		return nil, errors.New("不能重复领取任务")
	}

	task.Message = msg
	task.Price = detail.Price
	task.Final = detail.Final
	task.Quantity = detail.Quantity

	if ok := model.Insert(task); !ok {
		return nil, errors.New("领取任务失败")
	}

	return task, nil
}

func GetTaskByUserActivity(userID string, activityID int64) (*model.Task, error) {
	task := &model.Task{UserID: userID, ActivityID: activityID}
	if ok := model.Get(task); !ok {
		return nil, errors.New("未领取该活动任务")
	}
	task.DiscountNum = float64(task.Discount) / float64(Cardinal)
	return task, nil
}

func ListTasks(userID string, from, count int) (int64, interface{}) {
	task := &model.Task{UserID: userID}
	total, list := task.List(from, count)
	for i := range list {
		list[i].DiscountNum = float64(list[i].Discount) / float64(Cardinal)
	}
	return total, list
}

func GetTask(id int64) (*model.Task, error) {
	task := &model.Task{ID: id}
	if ok := model.Get(task); !ok {
		return nil, errors.New("指定任务不存在")
	}
	task.DiscountNum = float64(task.Discount) / float64(Cardinal)
	return task, nil
}

func DeleteTask(id int64) (*model.Task, error) {
	task := &model.Task{ID: id}
	if ok := model.Delete(task); !ok {
		return nil, errors.New("指定任务不存在")
	}
	return task, nil
}
