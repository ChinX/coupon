package model

import "time"

const (
	TaskDoing = iota
	TaskWaiting
	TaskDone
)

func init() {
	register(&OfficialAccount{}, &Admin{}, &User{}, &Activity{}, &Task{}, &Bargain{})
}

type OfficialAccount struct {
	ID        int64     `json:"id" xorm:"id notnull pk autoincr"`
	Account   string    `json:"-" xorm:"account varchar(64) notnull"`
	Nickname  string    `json:"nickname" xorm:"nickname varchar(64) notnull"`
	AvatarURL string    `json:"avatar_url" xorm:"avatar_url varchar(255) notnull"`
	QRCode    string    `json:"qr_code" xorm:"qr_code varchar(255) notnull"`
	CreatedAt time.Time `json:"-" xorm:"created"`
	DeletedAt time.Time `json:"-" xorm:"deleted"`
}

func (o *OfficialAccount) TableName() string {
	return "wxofficial"
}

type Admin struct {
	User     string `xorm:"user varchar(40) notnull pk"`
	Password string `xorm:"password varchar(128) notnull"`
	Salt     string `xorm:"salt varchar(8) notnull"`
}

func (a *Admin) TableName() string {
	return "admin"
}

type User struct {
	ID        string    `json:"openid" xorm:"id varchar(40) notnull pk"`
	Nickname  string    `json:"nickName"xorm:"nickname varchar(30) notnull"`
	AvatarURL string    `json:"avatarUrl"xorm:"avatar_url varchar(255)"`
	Gender    int       `json:"gender" xorm:"gender notnull default 0"`
	City      string    `json:"city" xorm:"city varchar(40)"`
	Province  string    `json:"province" xorm:"province varchar(40)"`
	Country   string    `json:"country" xorm:"country varchar(40)"`
	Language  string    `json:"language" xorm:"language varchar(20)"`
	CreatedAt time.Time `json:"-" xorm:"created"`
}

func (u *User) TableName() string {
	return "user"
}

type Activity struct {
	ID             int64     `json:"id" xorm:"id notnull pk autoincr"`
	Title          string    `json:"title" xorm:"title varchar(255) notnull"`
	Country        string    `json:"country" xorm:"country varchar(40)"`
	Province       string    `json:"province" xorm:"province varchar(40)"`
	City           string    `json:"city" xorm:"city varchar(40)"`
	DetailURL      string    `json:"detail_url" xorm:"detail_url varchar(255)"`
	AvatarURL      string    `json:"avatar_url" xorm:"avatar_url varchar(255)"`
	Description    string    `json:"description" xorm:"description varchar(255)"`
	PublicityIMG   string    `json:"publicity_img" xorm:"publicity_img varchar(255)"`
	CreatedAt      time.Time `json:"-" xorm:"created"`
	DeletedAt      time.Time `json:"-" xorm:"deleted"`
	Price          int       `json:"price" xorm:"price notnull default 0"`
	Final          int       `json:"final" xorm:"final notnull default 0"`
	Quantity       int       `json:"-" xorm:"quantity notnull default 0"`
	Total          int64     `json:"total" xorm:"Total notnull default 0"`
	Completed      int64     `json:"completed" xorm:"completed notnull default 0"`
	DailyTotal     int64     `json:"daily_total" xorm:"daily_total notnull default 0"`
	DailyCompleted int64     `json:"daily_completed" xorm:"daily_completed notnull default 0"`
	CouponStarted  time.Time `json:"-"  xorm:"coupon_started notnull"`
	CouponEnded    time.Time `json:"-" xorm:"coupon_ended notnull"`
	Expire         int64     `json:"expire" xorm:"expire notnull"`
}

func (a *Activity) TableName() string {
	return "activity"
}

type Task struct {
	ID            int64     `json:"id" xorm:"id notnull pk autoincr"`
	Message       string    `json:"message" xorm:"message varchar(255) notnull"`
	Price         int       `json:"price" xorm:"price notnull default 0"`
	Final         int       `json:"final" xorm:"final notnull default 0"`
	Discount      int       `json:"-" xorm:"discount notnull default 0"`
	DiscountNum   float64   `json:"discount" xorm:"-"`
	Quantity      int       `json:"-" xorm:"quantity notnull default 0"`
	Progress      int       `json:"progress" xorm:"progress notnull default 0"`
	UserID        string    `xorm:"user_id varchar(40) notnull unique(UQE_USER_ACTIVITY)"`
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

type Bargain struct {
	ID          int64     `json:"id" xorm:"id notnull pk autoincr"`
	Discount    int       `json:"-" xorm:"discount notnull default 0"`
	DiscountNum float64   `json:"discount" xorm:"-"`
	Message     string    `json:"message" xorm:"message varchar(255)"`
	UserID      string    `json:"user_id" xorm:"user_id varchar(40) notnull unique(UQE_USER_TASK)"`
	Nickname    string    `json:"nickName" xorm:"-"`
	AvatarURL   string    `json:"avatar_url" xorm:"-"`
	TaskID      int64     `json:"task_id" xorm:"task_id notnull unique(UQE_USER_TASK)"`
	CreatedAt   time.Time `json:"-" xorm:"created"`
	DeletedAt   time.Time `json:"-" xorm:"deleted"`
}

func (b *Bargain) TableName() string {
	return "bargain"
}
