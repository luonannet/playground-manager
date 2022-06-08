package util

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//app config
type Config struct {
	AppName       string         `json:"app_name"`
	AppModel      string         `json:"app_model"`
	AppHost       string         `json:"app_host"`
	AppPort       int            `json:"app_port"`
	Database      DatabaseConfig `json:"database"`
	K8sConfig     K8sConfig      `json:"k8s"`
	AuthingConfig AuthingConfig  `json:"authing"`
	JwtConfig     JwtConfig      `json:"jwt"`
	Statistic     Statistic      `json:"statistic"`
}

type K8sConfig struct {
	Namespace string `json:"namespace"`
	Image     string `json:"image"`
	FfileType string `json:"ffileType"`
}

//sql config
type DatabaseConfig struct {
	Driver   string `json:"driver"`
	DBUser   string `json:"db_user"`
	Password string `json:"password"`
	DBHost   string `json:"db_host"`
	DBPort   string `json:"db_port"`
	DbName   string `json:"db_name"`
	Chartset string `json:"charset"`
	ShowSql  bool   `json:"show_sql"`
}

//Authing Config
type AuthingConfig struct {
	UserPoolID  string `json:"userPoolID"`
	Secret      string `json:"secret"`
	AppID       string `json:"appID"`
	AppSecret   string `json:"appSecret"`
	RedirectURI string `json:"redirect_uri"`
}

//Jwt Jwt
type JwtConfig struct {
	Expire int    `json:"expire"`
	JwtKey string `json:"jwtKey"`
}

//Statistic function
type Statistic struct {
	Dir           string `json:"dir"`
	LogFile       string `json:"log_file"`
	LogFileSize   int64  `json:"log_file_size"`
	LogFileSuffix string `json:"log_file_suffix"`
}

func InitConfig(path string) {
	if dir, err := os.Getwd(); err == nil {
		if path == "" {
			dir = dir + "/conf/app.json"
		} else {
			dir = path
		}
		err = parseConfig(dir)
		if err != nil {
			Log.Errorf("load app.json failed, app must exit .please check app.json path:%s,and error:%s", dir, err)
			os.Exit(1)
			return
		}
	}

	if os.Getenv("GIN_MODE") != "" {
		cfg.AppModel = os.Getenv("GIN_MODE")
	}
	if cfg.AppModel == gin.DebugMode {
		//info level
		Log.SetLevel(logrus.InfoLevel)
	} else {
		//info level
		Log.SetLevel(logrus.ErrorLevel)
	}
	if os.Getenv("APP_PORT") != "" {
		cfg.AppPort, _ = strconv.Atoi(os.Getenv("APP_PORT"))
	}

	if os.Getenv("DB_USER") != "" {
		cfg.Database.DBUser = os.Getenv("DB_USER")
	}
	if os.Getenv("DB_PSWD") != "" {
		cfg.Database.Password = os.Getenv("DB_PSWD")
	}
	if os.Getenv("DB_HOST") != "" {
		cfg.Database.DBHost = os.Getenv("DB_HOST")
	}
	if os.Getenv("DB_NAME") != "" {
		cfg.Database.DbName = os.Getenv("DB_NAME")
	}

	if os.Getenv("AUTHING_APP_ID") != "" {
		cfg.AuthingConfig.AppID = os.Getenv("AUTHING_APP_ID")
	}
	if os.Getenv("AUTHING_APP_SECRET") != "" {
		cfg.AuthingConfig.AppSecret = os.Getenv("AUTHING_APP_SECRET")
	}
	if os.Getenv("AUTHING_SECRET") != "" {
		cfg.AuthingConfig.Secret = os.Getenv("AUTHING_SECRET")
	}
	if os.Getenv("AUTHING_USER_POOL_ID") != "" {
		cfg.AuthingConfig.UserPoolID = os.Getenv("AUTHING_USER_POOL_ID")
	}
	if os.Getenv("JWT_KEY") != "" {
		cfg.JwtConfig.JwtKey = os.Getenv("JWT_KEY")
	}
}

//external
func GetConfig() *Config {
	return cfg
}

//internal
var cfg *Config = nil

func parseConfig(path string) error {
	file, err := os.Open(path)
	if err != nil {
		Log.Errorf("read config file failed, please check path .  app exit now .")
		os.Exit(1)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)
	if err = decoder.Decode(&cfg); err != nil {
		return err
	}
	return nil
}
