package module

import (
	"errors"
	"math/rand"
	"time"

	"github.com/chinx/coupon/dao/mysql"
	"github.com/chinx/coupon/model"
)

const (
	min        = 1
	Cardinal   = 100
	TimeFormat = "2006年01月02日"

	TaskDoing = iota
	TaskWaiting
	TaskDone
)

func GetUserInfo(userID int64) (*model.User, bool) {
	user := &model.User{}
	return user, mysql.Get(user, "id=?", userID) == nil
}

func ListActivities(userID int64, from, count int) map[string]interface{} {
	activity := &model.Activity{}
	timeNow := time.Now()
	total, list := mysql.List(activity, from, count, "")
	account := &model.WXOfficial{}
	mysql.Get(account, "id > 0")
	account.AvatarURL = urlFormat(account.AvatarURL)
	account.QRCode = urlFormat(account.QRCode)

	l := len(list)
	for i := 0; i < l; i++ {
		item := list[i].(*model.Activity)
		item.DetailURL = urlFormat(item.DetailURL)
		item.PublicityIMG = urlFormat(item.PublicityIMG)
		item.AvatarURL = urlFormat(item.AvatarURL)
		item.Expire = int64(item.EndedAt.Sub(timeNow).Seconds())
		if item.Expire < 0 {
			item.DeletedAt = timeNow
			mysql.Update(item, "id=?", item.ID)
			preList := list[:i]
			if i+1 < l {
				preList = append(preList, list[i+1:]...)
			}
			list = preList
			i--
			l--
			total--
		}
	}

	task := &model.Task{}

	if userID > 0 {
		mysql.GetLast(task, "id", "user_id=?", userID)
	}

	return map[string]interface{}{
		"total":            total,
		"list":             list,
		"official_account": account,
		"activity_id":      task.ActivityID,
	}
}

func ListBargains(taskID int64, from, count int) interface{} {
	bargain := &model.Bargain{}
	total, list := mysql.List(bargain, from, count, "task_id=?", taskID)
	userMap := make(map[int64]*model.User)

	for i := range list {
		item := list[i].(*model.Bargain)
		user, ok := userMap[item.UserID]
		if !ok {
			user = &model.User{}
			mysql.Get(user, "id=?", item.UserID)
			userMap[item.UserID] = user
		}
		item.DiscountNum = float64(item.Discount) / float64(Cardinal)
		item.Nickname = user.Nickname
		item.AvatarURL = user.AvatarURL
	}
	return map[string]interface{}{
		"total": total,
		"list":  list,
	}
}

func UserTask(selfID int64, userID int64, activityID int64) (map[string]interface{}, error) {
	activity := &model.Activity{}
	err := mysql.Get(activity, "id=?", activityID)
	if err != nil {
		return nil, errors.New("指定任务不存在")
	}

	activity.DetailURL = urlFormat(activity.DetailURL)
	activity.PublicityIMG = urlFormat(activity.PublicityIMG)
	activity.AvatarURL = urlFormat(activity.AvatarURL)
	activity.Expire = int64(activity.EndedAt.Sub(time.Now()).Seconds())

	task := &model.Task{}
	var selfBargain *model.Bargain
	err = mysql.Get(task, "user_id=? and activity_id=?", userID, activityID)
	if err != nil {
		if selfID == userID {
			task.UserID = userID
			task.ActivityID = activityID
			task.Message = "这里有一个好玩的地方，大家帮我们拿门票啊"
			task.Price = activity.Price
			task.Final = activity.Final
			task.Quantity = activity.Quantity
			task.CouponStarted = activity.CouponStarted
			task.CouponEnded = activity.CouponEnded
			task.Discount += bargainPrice(task.Price-task.Final, task.Discount, task.Quantity, task.Progress)

			session := mysql.NewSession()
			defer session.Close()

			err := session.Begin()
			err = session.Insert(task)
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

			err = session.Insert(bargain)
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
		if selfID > 0 && mysql.Exist(&model.Bargain{}, "user_id=? and task_id=?", selfID, task.ID) {
			task.Bargained = 1
		}
	}

	userMap := make(map[int64]*model.User)
	user := &model.User{}
	mysql.Get(user, "id=?", task.UserID)
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

func CreateBargain(userID int64, taskID int64) (map[string]interface{}, error) {
	task := &model.Task{ID: taskID}
	err := mysql.Get(task, "id=?", taskID)
	if err != nil {
		return nil, errors.New("指定的任务不存在")
	}

	if task.Progress == task.Quantity {
		return nil, errors.New("您来晚已一步，对方任务已完成")
	}

	bargain := &model.Bargain{}

	err = mysql.Get(bargain, "user_id=? and task_id=?", userID, task.ID)
	if err == nil {
		return nil, errors.New("不能重复砍刀")
	} else if err != mysql.NoRecords {
		return nil, errors.New("砍刀失败：" + err.Error())
	}

	bargain.UserID = userID
	bargain.TaskID = task.ID
	bargain.Message = "轻松砍价到0元"
	task.Progress += 1

	count := bargainPrice(task.Price-task.Final, task.Discount, task.Quantity, task.Progress)
	task.Discount += count
	bargain.Discount = count

	if task.Progress == task.Quantity {
		task.Status = TaskWaiting
	}

	task.DiscountNum = float64(task.Discount) / float64(Cardinal)
	bargain.DiscountNum = float64(bargain.Discount) / float64(Cardinal)

	session := mysql.NewSession()
	defer session.Close()

	err = session.Begin()
	err = session.Update(task, "id=?", taskID)
	if err != nil {
		session.Rollback()
		return nil, err
	}

	err = session.Insert(bargain)
	if err != nil {
		session.Rollback()
		return nil, err
	}
	err = session.Commit()
	if err != nil {
		return nil, err
	}
	user := &model.User{}
	mysql.Get(user, "id=?", userID)
	bargain.Nickname = user.Nickname
	bargain.AvatarURL = user.AvatarURL

	task.Started = task.CouponStarted.Format(TimeFormat)
	task.Ended = task.CouponEnded.Format(TimeFormat)
	return map[string]interface{}{"task": task, "bargain": bargain}, nil
}

func CreateCash(userID int64, taskID int64) (*model.Task, error) {
	task := &model.Task{ID: taskID}
	err := mysql.Get(task, "user_id=?", userID)
	if err != nil {
		return nil, errors.New("指定的任务不存在")
	}

	if task.Status == TaskDone {
		return nil, errors.New("不能重复兑换")
	}

	if task.Status != TaskWaiting || task.Progress != task.Quantity {
		return nil, errors.New("兑换条件未达成")
	}

	task.Status = TaskDone
	err = mysql.Update(task, "id=?", taskID)
	if err != nil {
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
