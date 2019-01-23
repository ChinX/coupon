package api

type UserLogin struct {
	Code string `json:"code"`
}

type PageParams struct {
	From  int `json:"from"`
	Count int `json:"count"`
}

type ReplyResult struct {
	UserID     int64       `json:"user_id"`
	Status     int         `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	StatusCode int         `json:"-"`
	Error      error       `json:"-"`
}
