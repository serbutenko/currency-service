package exchangeratehost

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	apiKey string
}

func NewClient(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

type convertResponse struct {
	Info struct {
		Quote float64 `json:"quote"`
	} `json:"info"`
}

type listResponse struct {
	List map[string]string `json:"currencies"`
}

func (c *Client) FetchRate(ctx context.Context, from, to string) (float64, error) {
	url := fmt.Sprintf("https://api.exchangerate.host/convert?access_key=%s&from=%s&to=%s&amount=1",
		c.apiKey, from, to,
	)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data convertResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}
	return data.Info.Quote, nil
}

func (c *Client) FetchCurrencies(ctx context.Context) (map[string]string, error) {
	url := fmt.Sprintf("https://api.exchangerate.host/list?access_key=%s", c.apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data listResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}
	return data.List, nil
}
