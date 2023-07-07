package boot

import (
	"os"
	"time"

	"github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/natefinch/lumberjack"
	uberZap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"yasf.com/backend/playground/kratos-layout/internal/conf"
)

type BootLog struct {
	conf *conf.Bootstrap
}

func (b *BootLog) Run() log.Logger {
	zapLogger := b.Setting()
	logger := log.With(zap.NewLogger(zapLogger),
		"time", log.Timestamp("2006-01-02 15:04:05.000"),
		"caller", log.DefaultCaller,
		"service_id", b.conf.Server.Id,
		"service_name", b.conf.Server.Name,
		"service_version", b.conf.Server.Version,
		"trace_id", tracing.TraceID(),
		"span_id", tracing.SpanID())
	return logger
}

func (b *BootLog) Setting(opts ...uberZap.Option) *uberZap.Logger {
	var core zapcore.Core
	coreLevel, _ := zapcore.ParseLevel(b.conf.Logger.Level)
	core = zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
		}), // 编码器配置
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
			// zapcore.AddSync(b.writeSyncer())
		),
		uberZap.NewAtomicLevelAt(coreLevel), // 日志级别
	)
	return uberZap.New(core, opts...)
}

func (b *BootLog) writeSyncer() zapcore.WriteSyncer {
	if b.conf.Server.Environment == "local" {
		dateTime := time.Now().Format("2006-01-02-15")
		lumberJackLogger := &lumberjack.Logger{
			Filename:   b.conf.Logger.Path + dateTime + ".log",
			MaxSize:    int(b.conf.Logger.MaxSize), //默认1MB
			MaxBackups: int(b.conf.Logger.MaxBackups),
			MaxAge:     int(b.conf.Logger.MaxAge),
			Compress:   b.conf.Logger.Compress,
		}
		return zapcore.AddSync(lumberJackLogger)
	}
	return zapcore.AddSync(&lumberjack.Logger{})
}
