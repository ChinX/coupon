package router

import (
	"net/http"

	"github.com/chinx/coupon/handler"
	"github.com/naoina/denco"
)

func InitRouter() (http.Handler, error) {
	mux := denco.NewMux()
	return mux.Build([]denco.Handler{
		mux.POST("/v1/login", handler.UserLogin),

		mux.POST("/v1/activities", handler.CreateActivity),
		mux.GET("/v1/activities", handler.ListActivities),
		mux.PATCH("/v1/activities/:activity_id", handler.ModifyActivity),
		mux.GET("/v1/activities/:activity_id", handler.GetActivity),
		mux.DELETE("/v1/activities/:activity_id", handler.DeleteActivity),

		mux.POST("/v1/tasks", handler.CreateTask),
		mux.GET("/v1/tasks", handler.ListTasks),
		mux.GET("/v1/tasks/:task_id", handler.GetTask),
		mux.DELETE("/v1/tasks/:task_id", handler.DeleteTask),

		mux.POST("/v1/bargains", handler.CreateBargain),
		mux.GET("/v1/bargains", handler.ListBargains),

		mux.GET("/v1/coupons", handler.ListCoupons),
		mux.GET("/v1/coupons/:coupon", handler.GetCoupon),
		mux.PUT("/v1/coupons/:coupon", handler.ModifyCoupon),
	})
}
