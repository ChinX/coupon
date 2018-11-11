package module

import (
	"errors"
	"math/rand"
	"time"

	"github.com/chinx/coupon/model"
)

const (
	min      = 1
	Cardinal = 100
)

type TaskBargain struct {
	Task    *model.Task    `json:"task"`
	Bargain *model.Bargain `json:"bargain"`
}

func CreateBargain(userID string, taskID int64, msg string) (*TaskBargain, error) {
	task := &model.Task{ID: taskID}
	if ok := model.Get(task); !ok {
		return nil, errors.New("指定的任务不存在")
	}

	bargain := &model.Bargain{
		UserID: userID,
		TaskID: taskID,
	}

	if ok := model.Get(bargain); ok {
		return nil, errors.New("不能重复砍刀")
	}

	bargain.Message = msg

	isSelf := task.UserID == userID
	if !isSelf {
		task.Progress += 1
	}

	count := bargainPrice(task.Price-task.Final, task.Discount, task.Quantity, task.Progress)
	task.Discount += count
	bargain.Discount = count

	if task.Progress == task.Quantity {
		task.Status = model.TaskStatusWaiting
	}

	task.DiscountNum = float64(task.Discount) / float64(Cardinal)

	bargain.DiscountNum = float64(bargain.Discount) / float64(Cardinal)
	return &TaskBargain{Task: task, Bargain: bargain}, bargain.CreateByTask(task)
}

func ListBargains(taskID int64, from, count int) (int64, interface{}) {
	bargain := &model.Bargain{TaskID: taskID}
	total, list := bargain.List(from, count)
	for i := range list {
		list[i].DiscountNum = float64(list[i].Discount) / float64(Cardinal)
	}
	return total, list
}

func bargainPrice(total, discount, quantity, progress int) int {
	surplus := total*Cardinal - discount
	if quantity == progress || surplus <= min {
		return surplus
	}

	average := surplus / (quantity + 1 - progress)
	if average == min {
		return average
	}

	preset := (average >> 1) * 3
	safeNum := (surplus - preset) / (quantity + 1 - progress)
	rand.Seed(time.Now().UnixNano())
	return preset + rand.Intn(safeNum+min)
}
