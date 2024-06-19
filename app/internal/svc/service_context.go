package svc

import (
	"github.com/hu-1996/gormx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"reverse-turing/app/internal/config"
	"reverse-turing/app/internal/dal"
	"reverse-turing/app/internal/middleware"
	"reverse-turing/common/plugin"
)

type ServiceContext struct {
	Config            config.Config
	DB                *gorm.DB
	AppAuthMiddleware rest.Middleware
	CorsMiddleware    rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:            c,
		DB:                initMysql(c),
		AppAuthMiddleware: middleware.NewAppAuthMiddleware(c.Auth.PrivateKey).Handle,
		CorsMiddleware:    middleware.NewCorsMiddleware().Handle,
	}
}

func initMysql(config config.Config) *gorm.DB {
	var err error
	DB, err := gorm.Open(mysql.Open(config.Database.Source), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("mysql open error: %s", err)
	}

	if config.Database.AutoMigrate {
		if err = DB.AutoMigrate(dal.Game{}, dal.Agent{}, dal.Character{}, dal.GameCharacter{}, dal.Sentence{}); err != nil {
			log.Fatalf("mysql migrate error: %s", err)
		}
	}

	logx.Info("migrate success")

	DB.Use(&plugin.BasePlugin{})
	gormx.Init(DB)
	logx.Info("mysql init success")
	return DB
}
