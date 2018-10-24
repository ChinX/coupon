package dao

import "time"

type Bargain struct {
	ID       int64     `xorm:"id notnull pk autoincr"`
	Discount int64     `xorm:"discount notnull default 0"`
	Message  string    `xorm:"message"`
	UserID   int64     `xorm:"user_id notnull unique(UQE_USER_TASK)"`
	TaskID   int64     `xorm:"task_id notnull unique(UQE_USER_TASK)"`
	Created  time.Time `xorm:"created notnull"`
}

func (b *Bargain) TableName() string {
	return "bargain"
}

func init() {
	register(&Bargain{})
}
