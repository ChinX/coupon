package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-session/session"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

type pagination struct {
	Total int64       `json:"total"`
	List  interface{} `json:"list"`
}

type pager interface {
	List(from, counter int) (int64, interface{})
}

func pagedQuery(r *http.Request, p pager) ([]byte, error) {
	r.ParseForm()
	from, err := strconv.Atoi(r.Form.Get("from"))
	if err != nil || from < 0 {
		from = 0
	}

	count, err := strconv.Atoi(r.Form.Get("count"))
	if err != nil || count < 30 {
		count = 30
	}

	return pagedResult(p.List(from, count))
}

func pagedResult(total int64, list interface{}) ([]byte, error) {
	return json.Marshal(&pagination{
		Total: total,
		List:  list,
	})
}

func readBody(body io.Reader, v interface{}) error {
	byt, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(byt, v)
	if err != nil {
		return err
	}
	return nil
}

func readUserID(w http.ResponseWriter, r *http.Request) (string,error)  {
	store, err := session.Start(context.Background(), w, r)
	if err != nil {
		return "", err
	}

	userID, ok := store.Get("openid")
	if !ok {
		return "", errors.New("获取用户信息失败")
	}

	return userID.(string), nil
}
