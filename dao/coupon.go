package dao

import "time"

type Coupon struct {
	ID         int64     `xorm:"id notnull pk autoincr"`
	UserID     int64     `xorm:"user_id notnull unique(UQE_USER_ACTIVITY)"`
	ActivityID int64     `xorm:"activity_id notnull unique(UQE_USER_ACTIVITY)"`
	Created    time.Time `xorm:"created notnull"`
	Deleted    time.Time `xorm:"deleted"`
}

func (c *Coupon) TableName() string {
	return "coupon"
}

func init() {
	register(&Coupon{})
}
