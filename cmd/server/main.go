package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"os"

	"encoding/json"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	currencyv1 "example.com/currency-service/api/currency/v1"
)

var apiKey = os.Getenv("API_KEY")

type server struct {
	currencyv1.UnimplementedCurrencyServiceServer
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

func (s *server) GetRate(ctx context.Context, req *currencyv1.GetRateRequest) (*currencyv1.GetRateResponse, error) {
	from := req.GetFrom()
	to := req.GetTo()
	url := fmt.Sprintf(
		"https://api.exchangerate.host/convert?access_key=%s&from=%s&to=%s&amount=1",
		apiKey, from, to,
	)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var data ConvertResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		panic(err)
	}

	return &currencyv1.GetRateResponse{
		From:      req.GetFrom(),
		To:        req.GetTo(),
		Rate:      data.Info.Quote,
		Timestamp: timestamppb.New(time.Now().UTC()),
	}, nil
}

func (s *server) Convert(ctx context.Context, req *currencyv1.ConvertRequest) (*currencyv1.ConvertResponse, error) {
	from := req.GetFrom()
	to := req.GetTo()
	amount := req.GetAmount()

	url := fmt.Sprintf(
		"https://api.exchangerate.host/convert?access_key=%s&from=%s&to=%s&amount=%f",
		apiKey, from, to, amount,
	)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var data ConvertResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		panic(err)
	}

	return &currencyv1.ConvertResponse{
		From:      req.GetFrom(),
		To:        req.GetTo(),
		Amount:    req.GetAmount(),
		Converted: data.Result,
		Rate:      data.Info.Quote,
	}, nil
}

func (s *server) ListCurrencies(ctx context.Context, req *currencyv1.ListCurrenciesRequest) (*currencyv1.ListCurrenciesResponse, error) {
	url := fmt.Sprintf(
		"https://api.exchangerate.host/list?access_key=%s",
		apiKey,
	)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var data ListResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		panic(err)
	}

	return &currencyv1.ListCurrenciesResponse{
		Currencies: data.List,
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

	grpcServer := grpc.NewServer()
	currencyv1.RegisterCurrencyServiceServer(grpcServer, &server{})

	reflection.Register(grpcServer)

	fmt.Println("gRPC server is listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
