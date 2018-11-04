package model

import "time"

type User struct {
	ID        string    `json:"openid" xorm:"id varchar(40) notnull pk"`
	Nickname  string    `json:"nickName"xorm:"nickname varchar(30) notnull"`
	AvatarURL string    `json:"avatarUrl"xorm:"avatar_url varchar(255)"`
	Gender    int       `json:"gender" xorm:"gender notnull default 0"`
	City      string    `json:"city" xorm:"city varchar(40)"`
	Province  string    `json:"province" xorm:"province varchar(40)"`
	Country   string    `json:"country" xorm:"country varchar(40)"`
	Language  string    `json:"language" xorm:"language varchar(20)"`
	CreatedAt time.Time `json:"created" xorm:"created"`
}

func (u *User) TableName() string {
	return "user"
}

func init() {
	register(&User{})
}
