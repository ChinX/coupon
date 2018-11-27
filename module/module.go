package module

import (
	"errors"
	"math/rand"
	"time"

	"github.com/chinx/coupon/model"
)

const (
	min        = 1
	Cardinal   = 100
	TimeFormat = "2006年01月02日"
)

func GetUserInfo(userID string) (*model.User, bool) {
	user := &model.User{ID: userID}
	return user, model.Get(user)
}

func ListActivities(userID string, from, count int) map[string]interface{} {
	activity := &model.Activity{}
	total, list := activity.List(from, count)
	account := &model.OfficialAccount{}
	model.Get(account)
	account.AvatarURL = urlFormat(account.AvatarURL)
	account.QRCode = urlFormat(account.QRCode)

	for i := range list {
		list[i].DetailURL = urlFormat(list[i].DetailURL)
		list[i].PublicityIMG = urlFormat(list[i].PublicityIMG)
		list[i].AvatarURL = urlFormat(list[i].AvatarURL)
	}


	task := &model.Task{
		UserID: userID,
	}

	if userID != ""{
		task = task.Last()
	}

	return map[string]interface{}{
		"total":            total,
		"list":             list,
		"official_account": account,
		"activity_id":      task.ActivityID,
	}
}

func ListBargains(taskID int64, from, count int) interface{} {
	bargain := &model.Bargain{TaskID: taskID}
	total, list := bargain.List(from, count)
	userMap := make(map[string]*model.User)

	for key := range list {
		user, ok := userMap[list[key].UserID]
		if !ok {
			user = &model.User{ID: list[key].UserID}
			model.Get(user)
			userMap[list[key].UserID] = user
		}
		list[key].DiscountNum = float64(list[key].Discount) / float64(Cardinal)
		list[key].Nickname = user.Nickname
		list[key].AvatarURL = user.AvatarURL
	}
	return map[string]interface{}{
		"total": total,
		"list":  list,
	}
}

func UserTask(selfID string, userID string, activityID int64) (map[string]interface{}, error) {
	activity := &model.Activity{}
	if !model.Get(activity) {
		return nil, errors.New("指定任务不存在")
	}

	activity.DetailURL = urlFormat(activity.DetailURL)
	activity.PublicityIMG = urlFormat(activity.PublicityIMG)
	activity.AvatarURL = urlFormat(activity.AvatarURL)

	task := &model.Task{UserID: userID, ActivityID: activityID}
	var selfBargain *model.Bargain
	if !model.Get(task) {
		if selfID == userID {
			task.Message = "这里有一个好玩的地方，大家帮我们拿门票啊"
			task.Price = activity.Price
			task.Final = activity.Final
			task.Quantity = activity.Quantity
			task.CouponStarted = activity.CouponStarted
			task.CouponEnded = activity.CouponEnded
			task.Discount += bargainPrice(task.Price-task.Final, task.Discount, task.Quantity, task.Progress)

			session := model.NewSession()
			defer session.Close()

			err := session.Begin()
			_, err = session.Insert(task)
			if err != nil {
				session.Rollback()
				return nil, err
			}

			bargain := &model.Bargain{
				UserID:   userID,
				TaskID:   task.ID,
				Message:  "轻松砍价到0元",
				Discount: task.Discount,
			}

			_, err = session.Insert(bargain)
			if err != nil {
				session.Rollback()
				return nil, err
			}
			err = session.Commit()
			if err != nil {
				return nil, err
			}

			task.ShowDialog = 1
			bargain.DiscountNum = float64(bargain.Discount) / float64(Cardinal)
			selfBargain = bargain
		} else {
			return nil, errors.New("指定任务不存在")
		}
	} else {
		bargain := &model.Bargain{
			UserID: selfID,
			TaskID: task.ID,
		}
		if model.Get(bargain) {
			task.Bargained = 1
		}
	}

	userMap := make(map[string]*model.User)
	user := &model.User{ID: task.UserID}
	model.Get(user)
	userMap[task.UserID] = user
	task.DiscountNum = float64(task.Discount) / float64(Cardinal)
	task.Nickname = user.Nickname
	task.AvatarURL = user.AvatarURL
	task.Started = task.CouponStarted.Format(TimeFormat)
	task.Ended = task.CouponEnded.Format(TimeFormat)

	return map[string]interface{}{
		"activity": activity,
		"task":     task,
		"bargain":  selfBargain,
	}, nil
}

func CreateBargain(userID string, taskID int64) (map[string]interface{}, error) {
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

	bargain.Message = "轻松砍价到0元"
	task.Progress += 1

	count := bargainPrice(task.Final, task.Discount, task.Quantity, task.Progress)
	task.Discount += count
	bargain.Discount = count

	if task.Progress == task.Quantity {
		task.Status = model.TaskWaiting
	}

	task.DiscountNum = float64(task.Discount) / float64(Cardinal)
	bargain.DiscountNum = float64(bargain.Discount) / float64(Cardinal)

	session := model.NewSession()
	defer session.Close()

	err := session.Begin()
	_, err = session.ID(taskID).Update(task)
	if err != nil {
		session.Rollback()
		return nil, err
	}

	_, err = session.Insert(bargain)
	if err != nil {
		session.Rollback()
		return nil, err
	}
	err = session.Commit()
	if err != nil {
		return nil, err
	}
	user := &model.User{ID: userID}
	model.Get(user)
	bargain.Nickname = user.Nickname
	bargain.AvatarURL = user.AvatarURL

	task.Started = task.CouponStarted.Format(TimeFormat)
	task.Ended = task.CouponEnded.Format(TimeFormat)
	return map[string]interface{}{"task": task, "bargain": bargain}, nil
}

func CreateCash(userID string, taskID int64) (*model.Task, error) {
	task := &model.Task{ID: taskID}
	if ok := model.Get(task); !ok || task.UserID != userID {
		return nil, errors.New("指定的任务不存在")
	}

	if task.Status == model.TaskDone {
		return nil, errors.New("不能重复兑换")
	}

	if task.Status != model.TaskWaiting || task.Progress != task.Quantity {
		return nil, errors.New("兑换条件未达成")
	}

	task.Status = model.TaskDone
	if ok := model.Update(taskID, task); !ok {
		return nil, errors.New("兑票失败")
	}
	return task, nil
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
