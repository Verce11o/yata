package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Postgres   PostgresConfig `yaml:"postgres" env-required:"true"`
	HTTPServer HTTPServer     `yaml:"http_server" env-required:"true"`
	Services   Services       `yaml:"services" env-required:"true"`
	Mode       string         `yaml:"mode"`
}

type PostgresConfig struct {
	Host     string `yaml:"PostgresqlHost" env:"POSTGRESQL_HOST"`
	Port     string `yaml:"PostgresqlPort" env:"POSTGRESQL_PORT"`
	User     string `yaml:"PostgresqlUser" env:"POSTGRESQL_USERNAME"`
	Password string `yaml:"PostgresqlPassword" env:"POSTGRESQL_PASSWORD"`
	Name     string `yaml:"PostgresqlDbname" env:"POSTGRESQL_NAME"`
}

type HTTPServer struct {
	Port string `yaml:"port" env:"HTTPSERVER_PORT"`
}

type Services struct {
	Auth struct {
		Addr string `yaml:"addr"`
	} `yaml:"auth"`
}

func LoadConfig() *Config {
	var cfg Config

	if err := cleanenv.ReadConfig("config.yml", &cfg); err != nil {
		log.Fatalf("error while reading config file: %s", err)
	}
	return &cfg

}