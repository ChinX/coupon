package setting

import (
	"io/ioutil"

	"github.com/chinx/coupon/crypts"
	"gopkg.in/yaml.v2"
)

var options *Options

func Mode() *ModeConf {
	return options.Mode
}

func Listener() *ListenConf {
	return options.Listener
}

func Mysql() *MysqlConf {
	return options.Mysql
}

func Redis() *RedisConf {
	return options.Redis
}

func WeiXin() *WeiXinConf {
	return options.WeiXin
}

type Options struct {
	Mode     *ModeConf   `yaml:"mode"`
	Listener *ListenConf `yaml:"listener"`
	Mysql    *MysqlConf  `yaml:"mysql"`
	Redis    *RedisConf  `yaml:"redis"`
	WeiXin   *WeiXinConf `yaml:"wei_xin"`
}

type ModeConf struct {
	LogLevel    string `yaml:"log_level"`
	LogFile     string `yaml:"log_file"`
	EnablePprof int    `yaml:"enable_pprof"`
	MaxSize     int    `yaml:"max_size"`
	Interval    int    `yaml:"interval"`
}

type ListenConf struct {
	Protocol  string `yaml:"protocol"`
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	CrtFile   string `yaml:"crt_file"`
	KeyFile   string `yaml:"key_file"`
	StaticDir string `yaml:"static_dir"`
}

type MysqlConf struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type RedisConf struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	Database int    `yaml:"database"`
	Session  int    `yaml:"session"`
}

type WeiXinConf struct {
	AppID     string `yaml:"app_id"`
	AppSecret string `yaml:"app_secret"`
}

func LoadConfigFile(configFile, privateFile string) error {
	if err := LoadPrivateFile(privateFile); err != nil {
		return err
	}

	confData, err := LoadDecryptFile(configFile)
	if err != nil {
		return err
	}

	opts := &Options{}
	err = yaml.Unmarshal(confData, opts)
	if err != nil {
		return err
	}
	options = opts
	return nil
}

func LoadPrivateFile(privateFile string) error {
	if privateFile == "" {
		return nil
	}
	cryptData, err := ioutil.ReadFile(privateFile)
	if err != nil {
		return err
	}
	return crypts.DecryptPrivateKey(cryptData)
}

func LoadDecryptFile(filePath string) ([]byte, error) {
	cryptData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return crypts.AesDecrypt(cryptData)
}
