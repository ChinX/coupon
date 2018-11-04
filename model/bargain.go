package model

import "time"

type Bargain struct {
	ID        int64     `xorm:"id notnull pk autoincr"`
	Discount  int64     `xorm:"discount notnull default 0"`
	Message   string    `xorm:"message varchar(255)"`
	UserID    string    `xorm:"user_id varchar(40) notnull unique(UQE_USER_TASK)"`
	TaskID    int64     `xorm:"task_id notnull unique(UQE_USER_TASK)"`
	CreatedAt time.Time `xorm:"created"`
}

func (b *Bargain) TableName() string {
	return "bargain"
}

func init() {
	register(&Bargain{})
}
