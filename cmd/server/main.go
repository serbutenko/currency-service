package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"os"
	"strconv"

	"encoding/json"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"example.com/currency-service/internal/cache"

	currencyv1 "example.com/currency-service/api/currency/v1"
)

var apiKey = os.Getenv("API_KEY")

type server struct {
	currencyv1.UnimplementedCurrencyServiceServer
	redisCache *cache.RedisCache
}

type ConvertResponse struct {
	Success bool `json:"success"`
	Query   struct {
		From   string  `json:"from"`
		To     string  `json:"to"`
		Amount float64 `json:"amount"`
	} `json:"query"`
	Info struct {
		Timestamp int64   `json:"timestamp"`
		Quote     float64 `json:"quote"`
	} `json:"info"`
	Result float64 `json:"result"`
}

type ListResponse struct {
	Success bool              `json:"success"`
	List    map[string]string `json:"currencies"`
}

func (s *server) RateFromApi(ctx context.Context, from string, to string) (float64, error) {
	url := fmt.Sprintf(
		"https://api.exchangerate.host/convert?access_key=%s&from=%s&to=%s&amount=1",
		apiKey, from, to,
	)
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to get rate from API: %w", err)
	}
	defer resp.Body.Close()

	var data ConvertResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, fmt.Errorf("failed to convert reponse from API: %w", err)
	}
	rate := data.Info.Quote
	err = s.redisCache.Set(ctx, fmt.Sprintf("rate:%s:%s", from, to), strconv.FormatFloat(rate, 'f', 6, 64), time.Hour)
	if err != nil {
		return 0, fmt.Errorf("failed to set rate in Redis: %w", err)
	}
	return rate, nil
}

func (s *server) GetRate(ctx context.Context, req *currencyv1.GetRateRequest) (*currencyv1.GetRateResponse, error) {
	from := req.GetFrom()
	to := req.GetTo()
	var rate float64

	val, err := s.redisCache.Get(ctx, fmt.Sprintf("rate:%s:%s", from, to))
	if err != nil {
		rate, err = s.RateFromApi(ctx, from, to)
		if err != nil {
			return nil, err
		}
	} else {
		rate, _ = strconv.ParseFloat(val, 64)
	}

	return &currencyv1.GetRateResponse{
		From:      req.GetFrom(),
		To:        req.GetTo(),
		Rate:      rate,
		Timestamp: timestamppb.New(time.Now().UTC()),
	}, nil
}

func (s *server) Convert(ctx context.Context, req *currencyv1.ConvertRequest) (*currencyv1.ConvertResponse, error) {
	from := req.GetFrom()
	to := req.GetTo()
	amount := req.GetAmount()

	var rate float64

	val, err := s.redisCache.Get(ctx, fmt.Sprintf("rate:%s:%s", from, to))

	if err != nil {
		rate, err = s.RateFromApi(ctx, from, to)
		if err != nil {
			return nil, err
		}
	} else {
		rate, _ = strconv.ParseFloat(val, 64)
	}

	result := amount * rate

	return &currencyv1.ConvertResponse{
		From:      req.GetFrom(),
		To:        req.GetTo(),
		Amount:    req.GetAmount(),
		Converted: result,
		Rate:      rate,
	}, nil
}

func (s *server) ListCurrencies(ctx context.Context, req *currencyv1.ListCurrenciesRequest) (*currencyv1.ListCurrenciesResponse, error) {
	var currencies map[string]string
	val, err := s.redisCache.Get(ctx, "currencies:list")

	if err != nil {
		url := fmt.Sprintf(
			"https://api.exchangerate.host/list?access_key=%s",
			apiKey,
		)
		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to get rate from API: %w", err)
		}
		defer resp.Body.Close()

		var data ListResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, fmt.Errorf("failed to convert reponse from API: %w", err)
		}

		currencies = data.List

		jsonData, err := json.Marshal(currencies)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal currencies: %w", err)
		}

		err = s.redisCache.Set(ctx, "currencies:list", string(jsonData), time.Hour)
		if err != nil {
			return nil, fmt.Errorf("failed to set rate in Redis: %w", err)
		}
	} else {
		json.Unmarshal([]byte(val), &currencies)
	}

	return &currencyv1.ListCurrenciesResponse{
		Currencies: currencies,
	}, nil
}

func main() {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API key is missing! Please set API_KEY")
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	rdb := cache.NewRedisCache("localhost:6379", "", 0)
	s := &server{redisCache: rdb}

	grpcServer := grpc.NewServer()
	currencyv1.RegisterCurrencyServiceServer(grpcServer, s)

	reflection.Register(grpcServer)

	fmt.Println("gRPC server is listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
