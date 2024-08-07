package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"time"
)

var Config = new(Configure)

type Configure struct {
	*AppConfig   `mapstructure:"app"`
	*MySQLConfig `mapstructure:"db"`
	*RedisConfig `mapstructure:"redis"`
}

type AppConfig struct {
	Port     string `mapstructure:"port"`
	TimeZone string `mapstructure:"time_zone"`
}

type MySQLConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Pass     string `mapstructure:"pass"`
	Database string `mapstructure:"database"`
}

type RedisConfig struct {
	Port string `mapstructure:"port"`
}

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("..")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
		return
	}
	if err := viper.Unmarshal(Config); err != nil {
		fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("Configuration is modified...")
		if err := viper.Unmarshal(Config); err != nil {
			fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
		}
	})

	cstZone, err := time.LoadLocation(Config.AppConfig.TimeZone)
	if err != nil {
		panic(err)
	}
	time.Local = cstZone

	return
}
