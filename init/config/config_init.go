package config

import (
	"github.com/spf13/viper"
	"time"
)

func InitConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("../config") //开发环境
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(&Conf)
	if err != nil {
		return err
	}
	return nil
}

var Conf = Config{}

type Config struct {
	APP struct {
		Mode        string //主程序运行模式
		VerifyMode  string //邮箱验证模式	(smtp)-(mx)-(regex)
		ChunkSize   int
		RegexConfig struct {
			Goroutines int
			Regex      string
		}
		MxConfig struct {
			Mode       string //验证模式	(doh)-(default)
			Goroutines int
			Timeout    time.Duration
			DnsServer  string
		}
		SmtpConfig struct {
			Mode       string
			Goroutines int
			Timeout    time.Duration
			Port       string
			Keys       []string
		}
	}
	Data struct {
		Dir       bool
		Path      string
		WhiteList string
		BlackList string
	}
}
