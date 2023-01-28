package currency_api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type APIConfig struct {
	Link string
	Key  string
}

type Response struct {
	Data map[string]float64 `json:"data"`
}

func FetchCurrency(cfg APIConfig, baseCur, convCur string) (Response, error) {
	var resp Response
	url := fmt.Sprintf("%s?apikey=%s&base_currency=%s&currencies=%s", cfg.Link, cfg.Key, baseCur, convCur)
	req, err := http.Get(url)
	if err != nil {
		log.Printf("[While Fetch Currency Error occured]: %s\n", err.Error())
		return resp, err
	}

	defer req.Body.Close()
	b, err := io.ReadAll(req.Body)

	if err := json.Unmarshal(b, &resp); err != nil {
		return resp, err
	}
	return resp, err
}
