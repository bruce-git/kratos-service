package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"yasf.com/backend/playground/kratos-layout/internal/conf"
	"yasf.com/backend/playground/kratos-layout/pkg"

	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)

// Data .
type Data struct {
	db *gorm.DB
}

// NewData .
func NewData(conf *conf.Data, logger log.Logger) (*Data, func(), error) {
	defer func() {
		if err := recover(); err != nil {
			log.NewHelper(logger).Errorw("kind", "mysql", "error", err)
		}
	}()

	zapLog := pkg.NewZapGormLogV2(log.NewHelper(logger), "mysql")
	gormConfig := &gorm.Config{
		Logger: zapLog,
	}

	// mysql数据库连接
	db, err := gorm.Open(mysql.Open(conf.Database.Source), gormConfig)
	if err != nil {
		panic(err)
	}
	db.Use(&pkg.TracePlugin{})

	d := &Data{
		db: db,
	}

	return d, func() {
		log.NewHelper(logger).Info("message", "closing the data resources")
	}, nil
}
