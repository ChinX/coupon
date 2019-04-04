package httputils

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/chinx/cobweb"
	"github.com/chinx/coupon/api"
)

type pagination struct {
	Total int64       `json:"total"`
	List  interface{} `json:"list"`
}

func pageParams(r *http.Request) *api.PageParams {
	params := &api.PageParams{}
	params.Count, _ = strconv.Atoi(r.URL.Query().Get("count"))
	if params.Count < 0 || params.Count > 100 {
		params.Count = 30
	}

	pageNum, _ := strconv.Atoi(r.URL.Query().Get("from"))
	if pageNum < 1 {
		params.From = 0
	} else {
		params.From = (pageNum - 1) * params.Count
	}

	log.Println("pageParams", *params)
	return params
}

func urlParam(r *http.Request, key string) string {
	params := r.Context().Value(context.Background())
	if params == nil {
		return ""
	}
	log.Println("urlParam", params)
	return params.(cobweb.Params).Get(key)
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

	log.Println(string(byt))
	err = json.Unmarshal(byt, v)
	if err != nil {
		return err
	}
	return nil
}

func Reply(w http.ResponseWriter, status int, data interface{}, err error) {
	var result []byte
	switch t := data.(type) {
	case []byte:
		result = t
	case string:
		result = []byte(t)
	case error:
		result = []byte(t.Error())
	default:
		byteData, err := json.Marshal(data)
		if err != nil {
			status = http.StatusInternalServerError
			result = []byte(err.Error())
		} else {
			result = byteData
		}
	}

	log.Println(string(result))
	if status >= http.StatusBadRequest {
		if err != nil {
			log.Println(err)
		} else {
			log.Println(status, string(result))
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(result)
}