package api

type UserLogin struct {
	Code string `json:"code"`
}

type CommonResult struct {
	UserID  string `json:"user_id"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type Message struct {
	Message string `json:"message"`
}

type PageParams struct {
	From  int `json:"from"`
	Count int `json:"count"`
}
