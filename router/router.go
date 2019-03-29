package router

import (
	"github.com/chinx/cobweb"
	"github.com/chinx/coupon/handler"
	"github.com/urfave/negroni"
)

func SetRouters(n *negroni.Negroni) error {
	mux := cobweb.New()
	mux.Get("/", handler.HomeHandler)
	mux.Get("/S54107FZ3Q.txt", handler.VerificationHandler)

	mux.Post("/v1/user/login", handler.UserLogin)
	mux.Group("/v1", func() {
		mux.Post("/user/binding", handler.UserBinding)

		mux.Get("/activities", handler.ListActivities)
		mux.Post("/user/:user_id/:activity_id/task", handler.UserTask)
		mux.Post("/task/:task_id/bargains", handler.CreateBargain)
		mux.Get("/task/:task_id/bargains", handler.ListBargain)
		mux.Post("/task/:task_id/cash", handler.CreateCash)
	})

	mux.Group("/v1/source/", func() {
		mux.Get("/*filename", handler.StaticHandler)
	})

	mux.Group("/editor", func() {
		mux.Get("/*filename", handler.StaticHandler)
	})

	h, err := mux.Build()
	if err != nil {
		return err
	}

	n.UseHandler(h)
	return nil
}
