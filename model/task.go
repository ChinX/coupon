package model

import "time"

const (
	TaskStatusDoing = iota
	TaskStatusWaiting
	TaskStatusDone
)

type Task struct {
	ID          int64     `json:"id" xorm:"id notnull pk autoincr"`
	Message     string    `json:"message" xorm:"message varchar(255) notnull"`
	Price       int       `json:"price" xorm:"price notnull default 0"`
	Final       int       `json:"final" xorm:"final notnull default 0"`
	Discount    int       `json:"-" xorm:"discount notnull default 0"`
	DiscountNum float64   `json:"discount" xorm:"-"`
	Quantity    int       `json:"quantity" xorm:"quantity notnull default 0"`
	Progress    int       `json:"progress" xorm:"progress notnull default 0"`
	UserID      string    `json:"user_id" xorm:"user_id varchar(40) notnull unique(UQE_USER_ACTIVITY)"`
	ActivityID  int64     `json:"activity_id" xorm:"activity_id notnull unique(UQE_USER_ACTIVITY)"`
	CreatedAt   time.Time `json:"created" xorm:"created"`
	DeletedAt   time.Time `json:"-" xorm:"deleted"`
	Status      int       `json:"status" xorm:"status"` //0: 进行中, 1:待领取, 2:已领取
}

func (t *Task) TableName() string {
	return "task"
}

func init() {
	register(&Task{})
}

func (t *Task) List(from, count int) (int64, []*Task) {
	list := make([]*Task, 0)
	n, _ := engine.Count(t)
	if n == 0 {
		return 0, list
	}
	if from == 0 {
		engine.Desc("id").Limit(count, 0).Find(&list)
	} else {
		engine.Where("id < ?", from).Desc("id").Limit(count, 0).Find(&list)
	}
	return n, list
}
