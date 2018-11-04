package model

import "time"

type Coupon struct {
	ID         int64     `xorm:"id notnull pk autoincr"`
	UserID     string    `xorm:"user_id varchar(40) notnull unique(UQE_USER_ACTIVITY)"`
	ActivityID int64     `xorm:"activity_id notnull unique(UQE_USER_ACTIVITY)"`
	CreatedAt  time.Time `xorm:"created"`
	DeletedAt  time.Time `xorm:"deleted"`
}

func (c *Coupon) TableName() string {
	return "coupon"
}

func init() {
	register(&Coupon{})
}
