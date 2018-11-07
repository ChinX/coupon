package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func reply(w http.ResponseWriter, status int, data interface{}) {
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
	if status >= http.StatusBadRequest{
		log.Println(status, string(result))
	}
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
