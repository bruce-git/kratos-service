package main

import (
	"yasf.com/backend/playground/kratos-layout/internal/boot"
	"yasf.com/backend/playground/kratos-layout/internal/conf"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
)

func newApp(conf *conf.Server, logger log.Logger, hs *http.Server, gs *grpc.Server) *kratos.App {
	return kratos.New(
		kratos.ID(conf.Id),
		kratos.Name(conf.Name),
		kratos.Version(conf.Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			hs,
			gs,
		),
	)
}

func main() {
	bootstrapConf := boot.NewBootConf().Run()
	bootLogger := boot.NewBootLog(bootstrapConf).Run()
	boot.NewBootTrace(bootstrapConf).Run()

	app, cleanup, err := wireApp(bootstrapConf.Server, bootstrapConf.Data, bootLogger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
