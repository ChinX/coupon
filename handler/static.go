package handler

import (
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://www.windup.cn", http.StatusMovedPermanently)
}
