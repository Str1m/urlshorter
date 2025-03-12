package config

import (
	"github.com/spf13/viper"
	"log"
	"time"
)

type Server struct {
	Host        string        `mapstructure:"host"`
	Port        string        `mapstructure:"port"`
	Timeout     time.Duration `mapstructure:"timeout"`
	IdleTimeout time.Duration `mapstructure:"idle_timeout"`
}

type Storage struct {
	Host     string `mapstructure:"host"`
	DBName   string `mapstructure:"dbname"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Port     string `mapstructure:"port"`
}

type User struct {
	Name     string `mapstructure:"name" env-required:"true"`
	Password string `mapstructure:"password" env-required:"true"`
}

type Config struct {
	Env           string  `mapstructure:"env"`
	ServerConfig  Server  `mapstructure:"server"`
	StorageConfig Storage `mapstructure:"storage"`
	UserConfig    User    `mapstructure:"user"`
}

// TODO Добавить загрузку пути откуда-то (переменные окружения или параметры запуска)

func LoadConfig() Config {
	var cfg Config
	readFromConfigFile()
	readFromEnvFile()

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("fatal error config file: %v", err)
	}
	return cfg
}

func readFromConfigFile() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %v", err)
	}
}

func readFromEnvFile() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	if err := viper.MergeInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}
