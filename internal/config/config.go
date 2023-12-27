package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	Postgres   PostgresConfig `yaml:"postgres" env-required:"true"`
	HTTPServer HTTPServer     `yaml:"http_server" env-required:"true"`
	RabbitMQ   RabbitMQ       `yaml:"rabbitmq" env-required:"true"`
	Services   Services       `yaml:"services" env-required:"true"`
	App        App            `yaml:"app" env-required:"true"`
	Metrics    Metrics        `yaml:"metrics" env-required:"true"`
	Mode       string         `yaml:"mode"`
}

type App struct {
	JWT  JWTConfig `yaml:"jwt"`
	Port string    `yaml:"port"`
}

type JWTConfig struct {
	Secret        string `yaml:"secret"`
	TokenTTLHours int    `yaml:"token_ttl_hours"`
	Salt          string `yaml:"salt"`
}

type PostgresConfig struct {
	Host     string `yaml:"PostgresqlHost" env:"POSTGRESQL_HOST"`
	Port     string `yaml:"PostgresqlPort" env:"POSTGRESQL_PORT"`
	User     string `yaml:"PostgresqlUser" env:"POSTGRESQL_USERNAME"`
	Password string `yaml:"PostgresqlPassword" env:"POSTGRESQL_PASSWORD"`
	Name     string `yaml:"PostgresqlDbname" env:"POSTGRESQL_NAME"`
}

type RabbitMQ struct {
	Username     string `yaml:"username" env-required:"true"`
	Password     string `yaml:"password" env-required:"true"`
	Host         string `yaml:"host" env-required:"true"`
	Port         string `yaml:"port" env-required:"true"`
	ExchangeName string `yaml:"exchangeName" env-required:"true"`
	QueueName    string `yaml:"queueName" env-required:"true"`
	ConsumerTag  string `yaml:"consumerTag" env-required:"true"`
	BindingKey   string `yaml:"bindingKey" env-required:"true"`
}

type HTTPServer struct {
	Port string `yaml:"port" env:"HTTPSERVER_PORT"`
}

type Metrics struct {
	Jaeger struct {
		Endpoint string `yaml:"endpoint"`
	} `yaml:"jaeger"`
}

type Services struct {
	Auth struct {
		Addr string `yaml:"addr"`
	} `yaml:"auth"`
	Tweets struct {
		Addr string `yaml:"addr"`
	} `yaml:"tweets"`
	Comments struct {
		Addr string `yaml:"addr"`
	}
	Notifications struct {
		Addr string `yaml:"addr"`
	}
}

func LoadConfig() *Config {
	var cfg Config

	if err := cleanenv.ReadConfig("config.yml", &cfg); err != nil {
		log.Fatalf("error while reading config file: %s", err)
	}
	return &cfg

}
