package router

import (
	"net/http"

	"github.com/go-session/redis"
	"github.com/go-session/session"
	"github.com/chinx/cobweb"
	"github.com/chinx/coupon/handler"
)

func InitRouter() (http.Handler, error) {
	session.InitManager(
		session.SetStore(redis.NewRedisStore(&redis.Options{
			Addr: "127.0.0.1:6379",
			DB:   15,
		})),
	)

	mux := cobweb.New()
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
			mux.Get("/:activity_id", handler.GetActivity)
			mux.Delete("/:activity_id", handler.DeleteActivity)
		})

		mux.Group("/tasks", func() {
			mux.Post("/", handler.CreateTask)
			mux.Get("/", handler.ListTasks)
			mux.Get("/:task_id", handler.GetTask)
			mux.Delete("/:task_id", handler.DeleteTask)
		})

		mux.Group("/bargains", func() {
			mux.Post("/", handler.CreateBargain)
			mux.Get("/", handler.ListBargains)
		})

		mux.Group("/coupons", func() {
			mux.Get("/", handler.ListCoupons)
			mux.Get("/:coupon", handler.GetCoupon)
			mux.Put("/:coupon", handler.ModifyCoupon)
		})
	}, handler.CheckSession)

	return mux.Build()
}
