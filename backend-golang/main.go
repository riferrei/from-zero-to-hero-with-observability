package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"

	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"go.elastic.co/apm/module/apmgoredis"
	"go.elastic.co/apm/module/apmgorilla"
	"go.elastic.co/apm/module/apmot"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
)

const (
	eventDataset     string = "backend-golang.log"
	basePriceDefault string = "basePriceDefault"
)

var (
	redisClient apmgoredis.Client
	logger      *zap.Logger
	specialCars = []string{
		"ferrari", "mclaren",
		"lamborghini", "mercedes",
		"koenigsegg", "bugatti",
	}
)

func main() {

	// Initialize connection with Redis
	redisURL := os.Getenv("REDIS_URL")
	nativeClient := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	context := context.Background()
	redisClient = apmgoredis.Wrap(nativeClient).WithContext(context)
	defer redisClient.Close()

	// Create a logger based on Elastic ECS
	encoderConfig := ecszap.NewDefaultEncoderConfig()
	core := ecszap.NewCore(encoderConfig, os.Stdout, zap.DebugLevel)
	logger = zap.New(core, zap.AddCaller())

	router := mux.NewRouter()
	router.HandleFunc("/checkStatus", checkStatus)
	router.HandleFunc("/estimateValue", estimateValue)
	// Configure the Elastic APM and OpenTracing
	router.Use(apmgorilla.Middleware())
	opentracing.SetGlobalTracer(apmot.New())
	// Open the microservice for business...
	log.Fatal(http.ListenAndServe(":8888", router))

}

func estimateValue(writer http.ResponseWriter, request *http.Request) {

	queryParams := request.URL.Query()

	var brand string
	if len(queryParams["brand"]) == 1 {
		brand = queryParams["brand"][0]
	} else {
		err := errors.New("Parameter 'brand' was not provided")
		logger.Error("No brand", zap.Error(err),
			zap.String("event.dataset", eventDataset))
		panic(err)
	}
	var model string
	if len(queryParams["model"]) == 1 {
		model = queryParams["model"][0]
	} else {
		err := errors.New("Parameter 'model' was not provided")
		logger.Error("No model", zap.Error(err),
			zap.String("event.dataset", eventDataset))
		panic(err)
	}
	var year int
	if len(queryParams["year"]) == 1 {
		year, _ = strconv.Atoi(queryParams["year"][0])
	} else {
		err := errors.New("Parameter 'year' was not provided")
		logger.Error("No year", zap.Error(err),
			zap.String("event.dataset", eventDataset))
		panic(err)
	}

	var estimate Estimate

	calcEstSpan, ctx := opentracing.StartSpanFromContext(request.Context(), "calculateEstimate")
	calcEstSpan.SetTag("brand", brand)
	calcEstSpan.SetTag("model", model)
	calcEstSpan.SetTag("year", year)
	estimate = calculateEstimate(ctx, brand, model, year)
	calcEstSpan.Finish()

	bytes, err := json.Marshal(estimate)
	if err != nil {
		panic(err)
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.Write(bytes)

}

// Estimate struct
type Estimate struct {
	Estimate int    `json:"estimate"`
	Brand    string `json:"brand"`
	Model    string `json:"model"`
	Year     int    `json:"year"`
}

func calculateEstimate(ctx context.Context, brand string, model string, year int) Estimate {

	logger.Info("Value estimation for brand: "+brand,
		zap.String("event.dataset", eventDataset))

	estimate := Estimate{
		Brand: brand,
		Model: model,
		Year:  year,
	}

	basePriceSpan, ctx := opentracing.StartSpanFromContext(ctx, "getBasePrice")
	brand = strings.ToLower(brand)
	basePrice := getBasePrice(ctx, brand)
	basePriceSpan.Finish()

	markUp := int(((float64(5) * float64(basePrice)) / float64(100)))
	// Special cars have an additional markup based on market data
	for _, specialCar := range specialCars {
		if specialCar == brand {
			markUp += additionalMarkUp()
			break
		}
	}

	estimate.Estimate = basePrice + markUp
	return estimate

}

func getBasePrice(ctx context.Context, brand string) int {

	var basePrice int = 0

	tmpClient := redisClient.WithContext(ctx)
	value, err := tmpClient.Get(brand).Result()
	if err == redis.Nil {
		value, _ = tmpClient.Get(basePriceDefault).Result()
	}
	basePrice, _ = strconv.Atoi(value)
	return basePrice

}

func additionalMarkUp() int {
	logger.Debug("Use market data for additional mark up",
		zap.String("event.dataset", eventDataset))
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
