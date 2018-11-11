package router

import (
	"net/http"

	"github.com/chinx/cobweb"
	"github.com/chinx/coupon/handler"
)

func InitRouter() (http.Handler, error) {
	mux := cobweb.New()
	mux.Get("/", handler.HomeHandler)
	mux.Group("/v1/admin", func() {
		mux.Post("/login", handler.AdminLogin)
		mux.Post("/logout", handler.UserLogout)
	})

	mux.Group("/v1/user", func() {
		mux.Post("/login", handler.UserLogin)
		mux.Post("/binding", handler.UserBinding)
		mux.Post("/logout", handler.UserLogout)
	})

	mux.Group("/v1", func() {
		mux.Group("/activities", func() {
			mux.Post("/", handler.CreateActivity)
			mux.Get("/", handler.ListActivities)
			mux.Patch("/:activity_id", handler.ModifyActivity)
			mux.Delete("/:activity_id", handler.DeleteActivity)
			mux.Get("/:activity_id", handler.GetActivity)

			mux.Post("/:activity_id/tasks", handler.CreateTask)
			mux.Get("/:activity_id/tasks", handler.ActiveTask)

			mux.Post("/:activity_id/coupons", handler.CreateCoupon)
			mux.Get("/:activity_id/coupons", handler.ActiveCoupon)
		})

		mux.Group("/tasks", func() {
			mux.Get("/", handler.ListTasks)
			mux.Get("/:task_id", handler.GetTask)
			mux.Delete("/:task_id", handler.DeleteTask)

			mux.Post("/:task_id/bargains", handler.CreateBargain)
			mux.Get("/:task_id/bargains", handler.ListBargains)
		})

		mux.Group("/coupons", func() {
			mux.Get("/", handler.ListCoupons)
			mux.Get("/:coupon_id", handler.GetCoupon)
			mux.Delete("/:coupon_id", handler.DeleteCoupon)
		})
	})

	mux.Group("/editor", func() {
		mux.Get("/", handler.EditHandler)
		mux.Get("/*filename", handler.StaticHandler)
	})

	return mux.Build()
}
