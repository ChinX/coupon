package model

import "time"

type Task struct {
	ID         int64     `xorm:"id notnull pk autoincr"`
	Message    string    `xorm:"message varchar(255) notnull"`
	Price      int64     `xorm:"price notnull default 0"`
	Progress   int64     `xorm:"progress notnull default 0"`
	UserID     string    `xorm:"user_id varchar(40) notnull unique(UQE_USER_ACTIVITY)"`
	ActivityID int64     `xorm:"activity_id notnull unique(UQE_USER_ACTIVITY)"`
	CreatedAt  time.Time `json:"created" xorm:"created"`
	Completed  time.Time `xorm:"completed"`
	DeletedAt  time.Time `xorm:"deleted"`
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
