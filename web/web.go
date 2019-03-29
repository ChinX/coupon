package web

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/chinx/coupon/midware"
	"github.com/chinx/coupon/router"
	"github.com/urfave/negroni"
)

func RunWebServe(listenAddr, staticDir string, cfg *tls.Config) (err error) {
	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger(), negroni.NewStatic(http.Dir(staticDir)))
	midware.SetMidwares(n)

	err = router.SetRouters(n)
	if err != nil {
		log.Printf("set router error: %s", err)
	}

	return startListener(listenAddr, n, cfg)
}

func startListener(listenAddr string, handler http.Handler, cfg *tls.Config) (err error) {
	if cfg == nil {
		if err = http.ListenAndServe(listenAddr, handler); err != nil {
			log.Printf("Start http service error: %s", err)
		}
		return
	}

	server := &http.Server{Addr: listenAddr, TLSConfig: cfg, Handler: handler}
	if err = server.ListenAndServeTLS("", ""); err != nil {
		log.Printf("Start https service error: %s", err)
		return err
	}
	return
}
