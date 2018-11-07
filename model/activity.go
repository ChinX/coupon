package model

import "time"

type Activity struct {
	ID           int64     `json:"id" xorm:"id notnull pk autoincr"`
	Title        string    `json:"title" xorm:"title varchar(255) notnull"`
	Price        int64     `json:"price" xorm:"price notnull default 0"`
	Limit        int64     `json:"limit" xorm:"limit notnull default 0"`
	Count        int64     `json:"count" xorm:"count notnull default 0"`
	City         string    `json:"city" xorm:"city notnull varchar(40)"`
	Province     string    `json:"province" xorm:"province notnull varchar(40)"`
	DetailURL    string    `json:"detail_url" xorm:"detail_url varchar(255)"`
	PublicityIMG string    `json:"publicity_img" xorm:"publicity_img varchar(255)"`
	DailyLimit   int64     `json:"daily_limit" xorm:"daily_limit notnull default 0"`
	DailyCount   int64     `json:"daily_count" xorm:"daily_count notnull default 0"`
	CreatedAt    time.Time `json:"created" xorm:"created"`
	DeletedAt    time.Time `json:"-" xorm:"deleted"`
}

func (a *Activity) TableName() string {
	return "activity"
}

func init() {
	register(&Activity{})
}

func (a *Activity) List(from, count int) (int64, interface{}) {
	list := make([]*Activity, 0)
	n, _ := engine.Count(a)
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
