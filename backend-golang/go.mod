module backend-golang

go 1.15

require (
	github.com/go-redis/redis/v8 v8.8.2
	github.com/gomodule/redigo v1.8.4
	github.com/gorilla/mux v1.8.0
	github.com/kr/pretty v0.1.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.20.0
	go.opentelemetry.io/otel v0.20.0
	go.opentelemetry.io/otel/exporters/otlp v0.20.0
	go.opentelemetry.io/otel/metric v0.20.0
	go.opentelemetry.io/otel/sdk v0.20.0
	go.opentelemetry.io/otel/sdk/metric v0.20.0
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)
