package model

import (
	"github.com/chinx/coupon/dao/mysql"
)

func init() {
	mysql.Register(
		&WXOfficial{},
		&Admin{},
		&User{},
		&Activity{},
		&Task{},
		&Bargain{})
}
