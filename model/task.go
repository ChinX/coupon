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
