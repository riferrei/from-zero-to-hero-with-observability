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
	"time"

	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"go.elastic.co/apm/module/apmgorilla"
	"go.elastic.co/apm/module/apmot"
	"go.elastic.co/ecszap"
	"go.uber.org/zap"
)

const eventDataset string = "backend-golang.log"

var logger *zap.Logger

func main() {

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
	defer calcEstSpan.Finish()
	estimate = calculateEstimate(ctx, brand, model, year) // Function executed within the span

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

	var basePrice int = 0
	estimate := Estimate{
		Brand: brand,
		Model: model,
		Year:  year,
	}

	switch brand {
	case "Toyota":
		basePrice = 25000
	case "Lexus":
		basePrice = 35000
	case "Ford":
		basePrice = 20000
	case "Nissan":
		basePrice = 20000
	case "Tesla":
		basePrice = 60000
	case "Ferrari":
		basePrice = specialPriceCalculation()
	default:
		basePrice = 30000
	}

	estimate.Estimate = ((rand.Intn(50) - 5) + 1) * basePrice
	return estimate

}

func specialPriceCalculation() int {
	logger.Debug("Calculating special price for exotic car...",
		zap.String("event.dataset", eventDataset))
	time.Sleep(5 * time.Second)
	logger.Debug("Use market data with a fixed base price",
		zap.String("event.dataset", eventDataset))
	return rand.Intn(10) * 100000
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
