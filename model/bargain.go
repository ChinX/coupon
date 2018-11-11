package model

import (
	"time"
)

type Bargain struct {
	ID          int64     `json:"id" xorm:"id notnull pk autoincr"`
	Discount    int       `json:"-" xorm:"discount notnull default 0"`
	DiscountNum float64   `json:"discount" xorm:"-"`
	Message     string    `json:"message" xorm:"message varchar(255)"`
	UserID      string    `json:"user_id" xorm:"user_id varchar(40) notnull unique(UQE_USER_TASK)"`
	TaskID      int64     `json:"task_id" xorm:"task_id notnull unique(UQE_USER_TASK)"`
	CreatedAt   time.Time `json:"created_at" xorm:"created"`
}

func (b *Bargain) TableName() string {
	return "bargain"
}

func init() {
	register(&Bargain{})
}

func (b *Bargain) List(from, count int) (int64, []*Bargain) {
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

func (b *Bargain) CreateByTask(task *Task) error {
	session := NewSession()
	defer session.Close()

	err := session.Begin()

	_, err = session.Where("id = ?", task.ID).Update(task)
	if err != nil {
		session.Rollback()
		return err
	}

	_, err = session.Insert(b)
	if err != nil {
		session.Rollback()
		return err
	}
	return session.Commit()
}
