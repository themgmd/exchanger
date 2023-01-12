package currencies

import (
	"encoding/json"
	"fmt"
	"github.com/onemgvv/exchanger/internal/config"
	"io"
	"log"
	"net/http"
)

type Response struct {
	Data map[string]float64 `json:"data"`
}

func FetchCurrency(cfg *config.Config, baseCur, convCur string) (Response, error) {
	var resp Response
	url := fmt.Sprintf("%s?apikey=%s&base_currency=%s&currencies=%s", cfg.API.Link, cfg.API.Key, baseCur, convCur)
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
