package model

import "time"

type Bargain struct {
	ID          int64     `json:"id" xorm:"id notnull pk autoincr"`
	Discount    int       `json:"-" xorm:"discount notnull default 0"`
	DiscountNum float64   `json:"discount" xorm:"-"`
	Message     string    `json:"message" xorm:"message varchar(255)"`
	UserID      int64     `json:"user_id" xorm:"user_id notnull unique(UQE_USER_TASK)"`
	Nickname    string    `json:"nickName" xorm:"-"`
	AvatarURL   string    `json:"avatar_url" xorm:"-"`
	TaskID      int64     `json:"task_id" xorm:"task_id notnull unique(UQE_USER_TASK)"`
	CreatedAt   time.Time `json:"-" xorm:"created"`
	DeletedAt   time.Time `json:"-" xorm:"deleted"`
}

func (b *Bargain) TableName() string {
	return "bargain"
}
