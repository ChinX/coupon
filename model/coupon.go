package model

import "time"

type Coupon struct {
	ID         int64     `json:"id" xorm:"id notnull pk autoincr"`
	UserID     string    `json:"user_id" xorm:"user_id varchar(40) notnull unique(UQE_USER_ACTIVITY)"`
	ActivityID int64     `json:"activity_id" xorm:"activity_id notnull unique(UQE_USER_ACTIVITY)"`
	CreatedAt  time.Time `json:"created_at" xorm:"created"`
	DeletedAt  time.Time `json:"-" xorm:"deleted"`
}

func (c *Coupon) TableName() string {
	return "coupon"
}

func init() {
	register(&Coupon{})
}

func (c *Coupon) List(from, count int) (int64, interface{}) {
	list := make([]*Activity, 0)
	n, _ := engine.Count(c)
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
