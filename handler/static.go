package handler

import (
	"github.com/chinx/coupon/module"
	"html/template"
	"net/http"
)

var StaticDir = "./static"

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "http://www.windup.cn", http.StatusMovedPermanently)
}

func EditHandler(w http.ResponseWriter, r *http.Request) {
	result := checkAdmin(w, r)
	if result.Status != module.StatusLogin {
		return
	}
	t, err := template.ParseFiles(StaticDir + "/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	r.URL.Path = urlParam(r, "filename")
	http.FileServer(http.Dir(StaticDir)).ServeHTTP(w, r)
}