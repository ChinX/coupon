package main

import (
	"log"
	"os"
	"runtime"

	"github.com/chinx/coupon/cmd"
	"github.com/urfave/cli"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	runtime.GOMAXPROCS(runtime.NumCPU())
	cmdWeb := cmd.GetCmdWeb()

	app := cli.NewApp()

	app.Name = "Coupon"
	app.Usage = "mini-program service"
	app.Version = "v0.0.1"

	app.Commands = []cli.Command{
		cmdWeb,
	}

	app.Flags = append(app.Flags, []cli.Flag{}...)
	if err := app.Run(os.Args); err != nil{
		log.Println(err)
	}
}
