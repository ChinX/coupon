package model

import "time"

type WXOfficial struct {
	ID        int64     `json:"id" xorm:"id notnull pk autoincr"`
	Account   string    `json:"-" xorm:"account varchar(64) notnull"`
	Nickname  string    `json:"nickname" xorm:"nickname varchar(64) notnull"`
	AvatarURL string    `json:"avatar_url" xorm:"avatar_url varchar(255) notnull"`
	QRCode    string    `json:"qr_code" xorm:"qr_code varchar(255) notnull"`
	CreatedAt time.Time `json:"-" xorm:"created_at"`
	DeletedAt time.Time `json:"-" xorm:"created_at"`
}

func (o *WXOfficial) TableName() string {
	return "wx_official"
}
