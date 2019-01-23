package model

type Admin struct {
	ID       int64  `json:"id" xorm:"id notnull pk autoincr"`
	User     string `json:"user" xorm:"user varchar(40) notnull unique"`
	Password string `json:"password" xorm:"password varchar(128) notnull"`
	Salt     string `json:"salt" xorm:"salt varchar(8) notnull"`
}

func (a *Admin) TableName() string {
	return "admin"
}
