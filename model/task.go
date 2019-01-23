package model

import "time"

type Task struct {
	ID            int64     `json:"id" xorm:"id notnull pk autoincr"`
	Message       string    `json:"message" xorm:"message varchar(255) notnull"`
	Price         int       `json:"price" xorm:"price notnull default 0"`
	Final         int       `json:"final" xorm:"final notnull default 0"`
	Discount      int       `json:"-" xorm:"discount notnull default 0"`
	DiscountNum   float64   `json:"discount" xorm:"-"`
	Quantity      int       `json:"-" xorm:"quantity notnull default 0"`
	Progress      int       `json:"progress" xorm:"progress notnull default 0"`
	UserID        int64     `json:"user_id" xorm:"user_id notnull unique(UQE_USER_ACTIVITY)"`
	Nickname      string    `json:"nickName" xorm:"-"`
	AvatarURL     string    `json:"avatar_url" xorm:"-"`
	ActivityID    int64     `json:"activity_id" xorm:"activity_id notnull unique(UQE_USER_ACTIVITY)"`
	CouponStarted time.Time `json:"-"  xorm:"coupon_started notnull"`
	CouponEnded   time.Time `json:"-" xorm:"coupon_ended notnull"`
	Started       string    `json:"coupon_started" xorm:"-"`
	Ended         string    `json:"coupon_ended"  xorm:"-"`
	CreatedAt     time.Time `json:"-" xorm:"created"`
	DeletedAt     time.Time `json:"-" xorm:"deleted"`
	ShowDialog    int       `json:"show_dialog" xorm:"-"`
	Bargained     int       `json:"bargained" xorm:"-"`
	Status        int       `json:"status" xorm:"status"` //0: 进行中, 1:未兑换， 2:已兑换
}

func (t *Task) TableName() string {
	return "task"
}
