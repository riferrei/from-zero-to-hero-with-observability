package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/propagation"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
)

const (
	serviceName             = "backend-golang"
	basePriceDefault string = "base-price-default"
	exoticCars       string = "exotic-cars"
)

var (
	redisClient *redis.Client
)

func main() {

	// OpenTelemetry configuration

	ctx := context.Background()

	endpoint := os.Getenv("EXPORTER_ENDPOINT")
	driver := otlpgrpc.NewDriver(
		otlpgrpc.WithInsecure(),
		otlpgrpc.WithEndpoint(endpoint),
	)
	exporter, err := otlp.NewExporter(ctx, driver)
	if err != nil {
		log.Fatalf("%s: %v", "failed to create exporter", err)
	}

	res0urce, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String("1.0"),
		),
	)

	bsp := sdktrace.NewBatchSpanProcessor(exporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res0urce),
		sdktrace.WithSpanProcessor(bsp),
	)

	pusher := controller.New(
		processor.New(
			simple.NewWithExactDistribution(),
			exporter,
		),
		controller.WithExporter(exporter),
		controller.WithCollectPeriod(2*time.Second),
	)
	err = pusher.Start(ctx)
	if err != nil {
		log.Fatalf("%s: %v", "failed to start the controller", err)
	}
	defer func() { _ = pusher.Stop(ctx) }()

	otel.SetTracerProvider(tracerProvider)
	global.SetMeterProvider(pusher.MeterProvider())
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.Baggage{},
		propagation.TraceContext{}),
	)

	// Redis configuration

	redisURL := os.Getenv("REDIS_URL")
	redisClient = redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	defer redisClient.Close()

	// REST API configuration

	router := mux.NewRouter()
	router.HandleFunc("/checkStatus", checkStatus)
	router.HandleFunc("/estimateValue", estimateValue)
	router.Use(otelmux.Middleware(serviceName))
	log.Fatal(http.ListenAndServe(":8888", router))

}

// Estimate struct
type Estimate struct {
	Estimate int    `json:"estimate"`
	Brand    string `json:"brand"`
	Model    string `json:"model"`
	Year     int    `json:"year"`
}

func estimateValue(writer http.ResponseWriter, request *http.Request) {

	queryParams := request.URL.Query()

	var brand string
	if len(queryParams["brand"]) == 1 {
		brand = queryParams["brand"][0]
	}
	var model string
	if len(queryParams["model"]) == 1 {
		model = queryParams["model"][0]
	}
	var year int
	if len(queryParams["year"]) == 1 {
		year, _ = strconv.Atoi(queryParams["year"][0])
	}

	var estimate Estimate

	tracer := otel.Tracer(serviceName)
	ctx, span := tracer.Start(request.Context(), "calculateEstimate")
	span.SetAttributes(
		attribute.String("brand", brand),
		attribute.String("model", model),
		attribute.Int("year", year))
	estimate = calculateEstimate(ctx, brand, model, year)
	defer span.End()

	bytes, _ := json.Marshal(estimate)
	writer.Header().Add("Content-Type", "application/json")
	writer.Write(bytes)

}

func calculateEstimate(ctx context.Context, brand string, model string, year int) Estimate {

	tracer := otel.Tracer(serviceName)

	estimate := Estimate{
		Brand: brand,
		Model: model,
		Year:  year,
	}

	brand = strings.ToLower(brand)

	// Retrieve the base price for the car
	spanCtx, redisGetSpan := tracer.Start(ctx, "Redis GET")
	redisGetSpan.SetAttributes(attribute.String("brand", brand))
	basePrice, _ := redisClient.Get(spanCtx, brand).Int()
	redisGetSpan.End()

	if basePrice == 0 { // Retrieve the base price then...
		spanCtx, redisGetSpan = tracer.Start(ctx, "Redis GET")
		redisGetSpan.SetAttributes(attribute.String("brand", brand))
		basePrice, _ = redisClient.Get(spanCtx, basePriceDefault).Int()
		redisGetSpan.End()
	}

	// Calculate mark up of 5% on top of the base price
	markUp := int(((float64(5) * float64(basePrice)) / float64(100)))

	// Exotic cars have an additional markup
	spanCtx, redisSisMemberSpan := tracer.Start(ctx, "Redis SISMEMBER")
	redisSisMemberSpan.SetAttributes(attribute.String("brand", brand))
	isExotic := redisClient.SIsMember(spanCtx, exoticCars, brand).Val()
	redisSisMemberSpan.End()

	if isExotic {
		markUp += additionalMarkUp()
	}

	estimate.Estimate = basePrice + markUp
	return estimate

}

func additionalMarkUp() int {
	time.Sleep(5 * time.Second)
	return rand.Intn(3) * 10000
}

// Status struct
type Status struct {
	Status string `json:"status"`
}

func checkStatus(writer http.ResponseWriter, request *http.Request) {
	status := Status{"UP"}
	bytes, err := json.Marshal(status)
	if err != nil {
		panic(err)
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.Write(bytes)
}
