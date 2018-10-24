package dao

import "time"

type Task struct {
	ID         int64     `xorm:"id notnull pk autoincr"`
	Message    string    `xorm:"message notnull"`
	Price      int64     `xorm:"price notnull default 0"`
	Progress   int64     `xorm:"progress notnull default 0"`
	UserID     int64     `xorm:"user_id notnull unique(UQE_USER_ACTIVITY)"`
	ActivityID int64     `xorm:"activity_id notnull unique(UQE_USER_ACTIVITY)"`
	Created    time.Time `xorm:"created notnull"`
	Completed  time.Time `xorm:"completed"`
	Deleted    time.Time `xorm:"deleted"`
}

func (t *Task) TableName() string {
	return "task"
}

func init() {
	register(&Task{})
}
