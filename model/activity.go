package model

import "time"

type Activity struct {
	ID             int64     `json:"id" xorm:"id notnull pk autoincr"`
	Title          string    `json:"title" xorm:"title varchar(255) notnull"`
	Country        string    `json:"country" xorm:"country varchar(40)"`
	Province       string    `json:"province" xorm:"province varchar(40)"`
	City           string    `json:"city" xorm:"city varchar(40)"`
	Address        string    `json:"address" xorm:"address varchar(255)"`
	DetailURL      string    `json:"detail_url" xorm:"detail_url varchar(255)"`
	AvatarURL      string    `json:"avatar_url" xorm:"avatar_url varchar(255)"`
	Description    string    `json:"description" xorm:"description varchar(255)"`
	PublicityIMG   string    `json:"publicity_img" xorm:"publicity_img varchar(255)"`
	ShareIMG       string    `json:"share_img" xorm:"share_img varchar(255)"`
	CreatedAt      time.Time `json:"-"`
	DeletedAt      time.Time `json:"-"`
	EndedAt        time.Time `json:"-" xorm:"ended"`
	Price          int       `json:"price" xorm:"price notnull default 0"`
	Final          int       `json:"final" xorm:"final notnull default 0"`
	Quantity       int       `json:"-" xorm:"quantity notnull default 0"`
	Total          int64     `json:"total" xorm:"total notnull default 0"`
	Completed      int64     `json:"completed" xorm:"completed notnull default 0"`
	DailyTotal     int64     `json:"daily_total" xorm:"daily_total notnull default 0"`
	DailyCompleted int64     `json:"daily_completed" xorm:"daily_completed notnull default 0"`
	CouponStarted  time.Time `json:"-"  xorm:"coupon_started notnull"`
	CouponEnded    time.Time `json:"-" xorm:"coupon_ended notnull"`
	Expire         int64     `json:"expire" xorm:"-"`
}

func (a *Activity) TableName() string {
	return "activity"
}
