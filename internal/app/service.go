package app

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"example.com/currency-service/internal/adapters/exchangeratehost"
	"example.com/currency-service/internal/cache"
)

type CurrencyService struct {
	cache     cache.Cache
	apiClient *exchangeratehost.Client
}

func NewCurrencyService(cache cache.Cache, apiClient *exchangeratehost.Client) *CurrencyService {
	return &CurrencyService{cache: cache, apiClient: apiClient}
}

func (s *CurrencyService) GetRate(ctx context.Context, from, to string) (float64, error) {
	key := fmt.Sprintf("rate:%s:%s", from, to)
	val, err := s.cache.Get(ctx, key)
	if err == nil {
		if rate, err := strconv.ParseFloat(val, 64); err == nil {
			return rate, nil
		}
	}

	rate, err := s.apiClient.FetchRate(ctx, from, to)
	if err != nil {
		return 0, err
	}

	s.cache.Set(ctx, key, strconv.FormatFloat(rate, 'f', 6, 64), time.Hour)
	return rate, nil
}

func (s *CurrencyService) Convert(ctx context.Context, from, to string, amount float64) (float64, float64, error) {
	rate, err := s.GetRate(ctx, from, to)
	if err != nil {
		return 0, 0, err
	}
	return rate, amount * rate, nil
}

func (s *CurrencyService) ListCurrencies(ctx context.Context) (map[string]string, error) {
	key := "currencies:list"
	val, err := s.cache.Get(ctx, key)
	if err == nil {
		var currencies map[string]string
		if json.Unmarshal([]byte(val), &currencies) == nil {
			return currencies, nil
		}
	}

	currencies, err := s.apiClient.FetchCurrencies(ctx)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(currencies)
	s.cache.Set(ctx, key, string(data), time.Hour)

	return currencies, nil
}
