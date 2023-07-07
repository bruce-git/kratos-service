package boot

import (
	"os"

	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"yasf.com/backend/playground/kratos-layout/internal/conf"

	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

type BootTrace struct {
	conf *conf.Bootstrap
}

func (t *BootTrace) Run() {
	var tp *tracesdk.TracerProvider
	resource, err := t.resource()
	if err != nil {
		log.Errorw("trace_resource_error", err)
	}
	if t.conf.Trace.Batcher {
		exporter, err := t.exporter(t.conf.Trace.Exporter)
		if err != nil {
			log.Errorw("trace_exporter_error", err)
		}
		tp = tracesdk.NewTracerProvider(
			tracesdk.WithBatcher(exporter),
			tracesdk.WithResource(resource),
		)
	} else {
		tp = tracesdk.NewTracerProvider(
			tracesdk.WithResource(resource),
		)
	}
	otel.SetTracerProvider(tp)
}

func (t *BootTrace) exporter(types string) (tracesdk.SpanExporter, error) {
	var exporter tracesdk.SpanExporter
	var err error

	switch types {
	case "stdout":
		exporter, err = stdouttrace.New(
			stdouttrace.WithPrettyPrint())
	case "file":
		var osFile *os.File
		osFile, _ = os.Create(t.conf.Trace.TraceFilePath)
		exporter, err = stdouttrace.New(
			stdouttrace.WithWriter(osFile))
	case "jaeger":
		exporter, err = jaeger.New(
			jaeger.WithCollectorEndpoint(
				jaeger.WithEndpoint(t.conf.Trace.Endpoint),
			))
	}
	return exporter, err
}

func (t *BootTrace) resource() (*resource.Resource, error) {
	resource, error := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(t.conf.Server.Name),          //实例名称
			attribute.String("environment", t.conf.Server.Environment), // 相关环境
			attribute.String("ID", t.conf.Server.Version),              //版本
		),
	)
	return resource, error
}
