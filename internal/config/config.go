package config

import "github.com/ilyakaznacheev/cleanenv"

type ConfigDatabase struct {
	Port     string `env:"POSTGRES_PORT" env-default:"5432"`
	Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
	Name     string `env:"POSTGRES_DB" env-default:"postgres"`
	User     string `env:"POSTGRES_USER" env-default:"user"`
	Password string `env:"POSTGRES_PASSWORD"`
}

type ConfigGrpc struct {
	Port string `env:"PORT" env-default:"3333"`
	Host string `env:"HOST" env-default:"localhost"`
}

type ConfigSmtp struct {
	Port     string `env:"SMTP_PORT" env-default:"3333"`
	Host     string `env:"SMTP_HOST" env-default:"localhost"`
	Password string `env:"SMTP_PASSWORD" env-default:""`
	From     string `env:"SMTP_FROM"`
}

type MainConfig struct {
	GrpcConf     *ConfigGrpc
	DataBaseConf *ConfigDatabase
	SmtpConf     *ConfigSmtp
}

var cnfDB ConfigDatabase
var cnfGRPC ConfigGrpc
var cnfSMTP ConfigSmtp

func MustLoad() *MainConfig {
	err := cleanenv.ReadConfig(".env", &cnfDB)
	if err != nil {
		panic("Error while get config")
	}
	err = cleanenv.ReadConfig(".env", &cnfGRPC)
	if err != nil {
		panic("Error while get config")
	}
	err = cleanenv.ReadConfig(".env", &cnfSMTP)
	if err != nil {
		panic("Error while get config")
	}
	return &MainConfig{&cnfGRPC, &cnfDB, &cnfSMTP}
}
