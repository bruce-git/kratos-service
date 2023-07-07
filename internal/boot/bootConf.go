package boot

import (
	"flag"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/google/uuid"
	"yasf.com/backend/playground/kratos-layout/internal/conf"
)

var (
	flagConf string
)

func init() {
	//获取命令行参数
	flag.StringVar(&flagConf, "conf", "../../configs", "config path, eg: -conf config.yaml")
	flag.Parse()
}

type BootConf struct{}

func (c *BootConf) Run() *conf.Bootstrap {
	configParams := config.New(
		config.WithSource(
			file.NewSource(flagConf),
		),
	)
	//defer configParams.Close()
	if err := configParams.Load(); err != nil {
		panic(err)
	}
	/*	if err := configParams.Watch("nsq.enable", func(key string, value config.Value) {
			fmt.Printf("config changed: %s = %v\n", key, value)
			// 在这里写回调的逻辑
		}); err != nil {
			log.Error(err)
		}*/
	var bc conf.Bootstrap
	if err := configParams.Scan(&bc); err != nil {
		panic(err)
	}
	serverId, _ := uuid.NewUUID()
	bc.Server.Id = serverId.String()
	return &bc
}
