package grpcport

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	currencyv1 "example.com/currency-service/api/currency/v1"
	"example.com/currency-service/internal/app"
)

type Server struct {
	currencyv1.UnimplementedCurrencyServiceServer
	service *app.CurrencyService
}

func NewServer(service *app.CurrencyService) *Server {
	return &Server{service: service}
}

func (s *Server) GetRate(ctx context.Context, req *currencyv1.GetRateRequest) (*currencyv1.GetRateResponse, error) {
	rate, err := s.service.GetRate(ctx, req.GetFrom(), req.GetTo())
	if err != nil {
		return nil, err
	}

	return &currencyv1.GetRateResponse{
		From:      req.GetFrom(),
		To:        req.GetTo(),
		Rate:      rate,
		Timestamp: timestamppb.New(time.Now().UTC()),
	}, nil
}

func (s *Server) Convert(ctx context.Context, req *currencyv1.ConvertRequest) (*currencyv1.ConvertResponse, error) {
	rate, result, err := s.service.Convert(ctx, req.GetFrom(), req.GetTo(), req.GetAmount())
	if err != nil {
		return nil, err
	}

	return &currencyv1.ConvertResponse{
		From:      req.GetFrom(),
		To:        req.GetTo(),
		Amount:    req.GetAmount(),
		Converted: result,
		Rate:      rate,
	}, nil
}

func (s *Server) ListCurrencies(ctx context.Context, req *currencyv1.ListCurrenciesRequest) (*currencyv1.ListCurrenciesResponse, error) {
	currencies, err := s.service.ListCurrencies(ctx)
	if err != nil {
		return nil, err
	}

	return &currencyv1.ListCurrenciesResponse{
		Currencies: currencies,
	}, nil
}
