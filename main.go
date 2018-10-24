package main

import (
	"fmt"
	"github.com/chinx/coupon/dao"
	"github.com/chinx/coupon/router"
	"net/http"
	"os"
)

func main() {
	err := dao.InitORM("mysql", "root:vessel@(127.0.0.1:3306)/vessel?charset=utf8")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	handler,err := router.InitRouter()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	fmt.Println(http.ListenAndServe(":8088", handler))
}
