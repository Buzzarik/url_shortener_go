package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ENV      string         `yaml:"env"`
	Postgres ConfigPostgres `yaml:"postgres"`
	Server   ConfigServer   `yaml:"server_url_shortener"`
}

type ConfigPostgres struct {
	Driver       string 		`yaml:"driver"`
	Host         string 		`yaml:"host"`
	Port         int    		`yaml:"port"`
	Username     string 		`yaml:"username"`
	Password     string 		`yaml:"password"`
	Sslmode      string 		`yaml:"sslmode"`
	DbName       string 		`yaml:"db_name"`
	MaxOpenConns int    		`yaml:"max_open_conns"`
	MaxIdleConns int    		`yaml:"max_idle_conns"`
	MaxIdleTime  time.Duration 		`yaml:"max_idle_time"`
};

type ConfigServer struct {
	Host 			string 		`yaml:"host"`
	Port 			int 		`yaml:"port"`
	Timeout 		time.Duration  	`yaml:"timeout"`
	IdleTimeout 	time.Duration 	`yaml:"idle_timeout"`
	IdApi 			int 		`yaml:"ad_api"`
};

func LoadConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH");
	if configPath == ""{
		log.Fatal("CONFIG_PATH is not set");
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err){
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var config Config;

	if err := cleanenv.ReadConfig(configPath, &config); err != nil{
		log.Fatalf("cannot read config: %s", err)
	}

	return &config;
}