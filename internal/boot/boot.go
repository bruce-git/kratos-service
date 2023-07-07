package boot

import (
	"yasf.com/backend/playground/kratos-layout/internal/conf"
)

type Boot interface {
	Run()
	Setting()
}

func NewBootConf() *BootConf {
	return &BootConf{}
}

func NewBootLog(conf *conf.Bootstrap) *BootLog {
	return &BootLog{
		conf,
	}
}

func NewBootTrace(conf *conf.Bootstrap) *BootTrace {
	return &BootTrace{
		conf,
	}
}
