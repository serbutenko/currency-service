package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"example.com/currency-service/internal/adapters/exchangeratehost"
	"example.com/currency-service/internal/app"
	"example.com/currency-service/internal/cache"
	"example.com/currency-service/internal/config"
	"example.com/currency-service/internal/ports/grpc"

	currencyv1 "example.com/currency-service/api/currency/v1"
)

func main() {
	cfg := config.Load()

	rdb := cache.NewRedisCache(cfg.RedisAddr, "", 0)
	apiClient := exchangeratehost.NewClient(cfg.ApiKey)

	service := app.NewCurrencyService(rdb, apiClient)
	grpcServer := grpc.NewServer()

	currencyv1.RegisterCurrencyServiceServer(
		grpcServer,
		grpcport.NewServer(service),
	)

	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	fmt.Printf("gRPC server is listening on %s\n", cfg.GRPCAddr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
