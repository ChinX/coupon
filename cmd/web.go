package cmd

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/chinx/coupon/dao/mysql"
	"github.com/chinx/coupon/setting"
	"github.com/chinx/coupon/web"
	"github.com/go-redis/redis"
	sedis "github.com/go-session/redis"
	"github.com/go-session/session"
	"github.com/urfave/cli"
)

var (
	emptyConfig = "%s config is empty"

	modeFlag = cli.StringFlag{
		Name:  "mode",
		Value: "dev",
		Usage: "runtime mode; options: 'dev','prod'",
	}

	privateFlag = cli.StringFlag{
		Name:  "private",
		Value: "",
		Usage: "the file path of private key; empty means no encryption",
	}

	confFlag = cli.StringFlag{
		Name:  "conf",
		Value: "./conf/config.yaml",
		Usage: "configuration file path; defaults to './conf/config.yaml'",
	}
)

// GetCmdWeb get a client command
func GetCmdWeb() cli.Command {
	return cli.Command{
		Name:        "run",
		Usage:       "run coupon web service",
		Description: "coupon is a mini-program service.",
		Action:      runWeb,
		Flags:       []cli.Flag{modeFlag, privateFlag, confFlag},
	}
}

func runWeb(c *cli.Context) (err error) {
	if err = parseConf(c); err != nil {
		return
	}

	if err = initMysql(setting.Mysql()); err != nil {
		return
	}

	if err = initSession(setting.Redis()); err != nil {
		return
	}

	if err = initCache(setting.Redis()); err != nil {
		return
	}

	if err = initListener(setting.Listener()); err != nil {
		return
	}

	return
}

func initMysql(conf *setting.MysqlConf) (err error) {
	if conf == nil {
		err = fmt.Errorf(emptyConfig, "mysql")
		log.Println(err)
		return
	}
	err = mysql.InitORM("mysql", fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8",
		conf.User, conf.Password, conf.Host, conf.Port, conf.Database))
	if err != nil {
		log.Println(err)
	}
	return
}

func initSession(conf *setting.RedisConf) (err error) {
	if conf == nil {
		err = fmt.Errorf(emptyConfig, "session")
		log.Println(err)
		return
	}
	session.InitManager(
		session.SetStore(sedis.NewRedisStore(&sedis.Options{
			Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
			DB:       conf.Session,
			Password: conf.Password,
		})),
	)
	return
}

func initCache(conf *setting.RedisConf) (err error) {
	if conf == nil {
		err = fmt.Errorf(emptyConfig, "cache")
		log.Println(err)
		return
	}
	redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
		DB:       conf.Database,
		Password: conf.Password,
	})
	return
}

func initListener(conf *setting.ListenConf) (err error) {
	if conf == nil {
		err = fmt.Errorf(emptyConfig, "listen")
		log.Println(err)
		return
	}

	listenAddr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	var cfg *tls.Config
	switch conf.Protocol {
	case "http":
	case "https":
		cfg, err = loadTlsConfig(conf)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("%s not supported", conf.Protocol)
	}
	return web.RunWebServe(listenAddr, conf.StaticDir, cfg)
}

func loadTlsConfig(conf *setting.ListenConf) (*tls.Config, error) {
	certData, err := ioutil.ReadFile(conf.CrtFile)
	if err != nil {
		log.Printf("read https crt file error: %s", err)
		return nil, err
	}

	keyData, err := setting.LoadDecryptFile(conf.KeyFile)
	if err != nil {
		log.Printf("read https key file error: %s", err)
		return nil, err
	}

	certificate, err := tls.X509KeyPair(certData, keyData)
	if err != nil {
		log.Printf("create tls key pair error: %s", err)
		return nil, err
	}

	cfg := &tls.Config{
		MinVersion:               tls.VersionTLS10,
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
	return cfg, nil
}

func parseConf(c *cli.Context) (err error) {
	confFile := c.String(confFlag.Name)
	privateFile := c.String(privateFlag.Name)
	return setting.LoadConfigFile(confFile, privateFile)
}
