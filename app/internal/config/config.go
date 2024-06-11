package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Database struct {
		Source      string
		AutoMigrate bool
	}
	Auth struct {
		PrivateKey string
	}
}
