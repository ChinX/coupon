package model

import "time"

type Bargain struct {
	ID        int64     `xorm:"id notnull pk autoincr"`
	Discount  float64   `xorm:"discount notnull default 0"`
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

func (b *Bargain) List(from, count int) (int64, interface{}) {
	list := make([]*Bargain, 0)
	n, _ := engine.Count(b)
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
