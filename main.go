package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/chinx/coupon/model"
	"github.com/chinx/coupon/module"
	"github.com/chinx/coupon/router"
	"github.com/chinx/coupon/setting"
	"github.com/go-session/redis"
	"github.com/go-session/session"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	opt, err := setting.LoadConfigFile("./cert/coupon_private.key", "./conf/windup.toml")
	//opt, err := setting.LoadConfigFile("./cert/coupon_private.key", "./conf/windup.conf")
	if err != nil {
		log.Fatal(err)
	}

	err = model.InitORM("mysql",
		fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8",
			opt.Mysql.User, opt.Mysql.Password,
			opt.Mysql.Server, opt.Mysql.Port,
			opt.Mysql.Database))

	if err != nil {
		log.Fatal(err)
	}
	session.InitManager(
		session.SetStore(redis.NewRedisStore(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", opt.Redis.Server, opt.Redis.Port),
			DB:       opt.Redis.Session,
			Password: opt.Redis.Password,
		})),
	)
	module.AppID = opt.Weixin.AppID
	module.AppSecret = opt.Weixin.AppSecret

	serveHandler, err := router.InitRouter()
	if err != nil {
		log.Fatal(err)
	}

	certData, err := ioutil.ReadFile(opt.CrtFile)
	if err != nil {
		log.Fatal(err)
	}

	keyData, err := setting.LoadDecryptFile(opt.KeyFile)
	if err != nil {
		log.Fatal(err)
	}

	certificate, err := tls.X509KeyPair(certData, keyData)
	if err != nil {
		log.Fatal(err)
	}

	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		},
		Certificates: []tls.Certificate{certificate},
	}
	srv := &http.Server{
		Addr:         "0.0.0.0:" + strconv.Itoa(opt.HttpsPort),
		Handler:      serveHandler,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	err = srv.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatal(err)
	}
}
