package router

import (
	"net/http"

	"github.com/chinx/cobweb"
	"github.com/chinx/coupon/handler"
)

func InitRouter() (http.Handler, error) {
	mux := cobweb.New()
	mux.Get("/", handler.HomeHandler)

	mux.Group("/v1", func() {
		mux.Post("/user/login", handler.UserLogin)
		mux.Post("/user/binding", handler.UserBinding)

		mux.Get("/activities", handler.ListActivities)
		mux.Get("/user/:user_id/:activity_id/task", handler.UserTask)
		mux.Post("/task/:task_id/bargains", handler.CreateBargain)
		mux.Get("/task/:task_id/bargains", handler.ListBargain)
		mux.Post("/task/:task_id/cash", handler.CreateCash)
	})
	return mux.Build()
}
