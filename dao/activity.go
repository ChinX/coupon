package dao

import "time"

type Activity struct {
	ID           int64     `xorm:"id notnull pk autoincr"`
	Title        string    `xorm:"title notnull"`
	Price        int64     `xorm:"price notnull default 0"`
	Limit        int64     `xorm:"limit notnull default 0"`
	Count        int64     `xorm:"count notnull default 0"`
	DetailURL    string    `xorm:"detail_url"`
	PublicityIMG string    `xorm:"publicity_img"`
	DailyLimit   int64     `xorm:"daily_limit notnull default 0"`
	DailyCount   int64     `xorm:"daily_count notnull default 0"`
	Created      time.Time `xorm:"created notnull"`
	Deleted      time.Time `xorm:"deleted"`
}

func (a *Activity) TableName() string {
	return "activity"
}

func init() {
	register(&Activity{})
}
