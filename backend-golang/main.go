package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"go.elastic.co/apm/module/apmgorilla"
	"go.elastic.co/apm/module/apmot"
	"go.elastic.co/apm/module/apmredigo"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
)

const (
	eventDataset     string = "backend-golang.log"
	basePriceDefault string = "base-price-default"
	exoticCars       string = "exotic-cars"
)

var (
	redisPool *redis.Pool
	logger    *zap.Logger
)

func main() {

	// Create a Redis connection pool
	redisURL := os.Getenv("REDIS_URL")
	redisPool = &redis.Pool{
		MaxIdle:     5,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", redisURL)
		},
	}

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
	} else {
		err := errors.New("Parameter 'brand' was not provided")
		logger.Error("No brand", zap.Error(err),
			zap.String("event.dataset", eventDataset))
	}
	var model string
	if len(queryParams["model"]) == 1 {
		model = queryParams["model"][0]
	} else {
		err := errors.New("Parameter 'model' was not provided")
		logger.Error("No model", zap.Error(err),
			zap.String("event.dataset", eventDataset))
	}
	var year int
	if len(queryParams["year"]) == 1 {
		year, _ = strconv.Atoi(queryParams["year"][0])
	} else {
		err := errors.New("Parameter 'year' was not provided")
		logger.Error("No year", zap.Error(err),
			zap.String("event.dataset", eventDataset))
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
		logger.Error("Error marshaling JSON response", zap.Error(err),
			zap.String("event.dataset", eventDataset))
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.Write(bytes)

}

func calculateEstimate(ctx context.Context, brand string, model string, year int) Estimate {

	logger.Info("Value estimation for brand: "+brand,
		zap.String("event.dataset", eventDataset))

	estimate := Estimate{
		Brand: brand,
		Model: model,
		Year:  year,
	}

	brand = strings.ToLower(brand)

	// Retrieve the base price for the car
	redisConn := apmredigo.Wrap(redisPool.Get()).WithContext(ctx)
	defer redisConn.Close()
	basePrice, err := redis.Int(redisConn.Do("GET", brand))
	if err != nil {
		logger.Error(fmt.Sprintf("Error getting base price for '%s'", brand),
			zap.Error(err), zap.String("event.dataset", eventDataset))
	}
	if basePrice == 0 {
		basePrice, err = redis.Int(redisConn.Do("GET", basePriceDefault))
		if err != nil {
			logger.Error("Error getting base price default", zap.Error(err),
				zap.String("event.dataset", eventDataset))
		}
	}

	// Calculate mark up of 5% on top of the base price
	markUp := int(((float64(5) * float64(basePrice)) / float64(100)))

	// Exotic cars have an additional markup
	isExotic, err := redis.Bool(redisConn.Do("SISMEMBER", exoticCars, brand))
	if err != nil {
		logger.Error(fmt.Sprintf("Error checking if '%s' is exotic", brand),
			zap.Error(err), zap.String("event.dataset", eventDataset))
	}
	if isExotic {
		markUp += additionalMarkUp()
	}

	estimate.Estimate = basePrice + markUp
	return estimate

}

func additionalMarkUp() int {
	logger.Debug("Waiting for the market data...",
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
