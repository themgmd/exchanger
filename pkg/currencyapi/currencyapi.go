package currencyapi

import (
	"encoding/json"
	"exchanger/pkg/errors"
	"fmt"
	"io"
	"net/http"
)

type Options struct {
	Link string
	Key  string
}

type Option func(*Options) error

func Link(link string) Option {
	return func(options *Options) error {
		options.Link = link
		return nil
	}
}

func Key(key string) Option {
	return func(options *Options) error {
		options.Key = key
		return nil
	}
}

type Response struct {
	Data map[string]float64 `json:"data"`
}

type CurrencyApi struct {
	options Options
}

func New(options ...Option) *CurrencyApi {
	var option Options

	for i := range options {
		if err := options[i](&option); err != nil {
			return &CurrencyApi{}
		}
	}

	return &CurrencyApi{
		options: option,
	}
}

func (ca CurrencyApi) Fetch(baseCur, convCur string) (Response, error) {
	var resp Response
	url := fmt.Sprintf("%s?apikey=%s&base_currency=%s&currencies=%s",
		ca.options.Link, ca.options.Key, baseCur, convCur)

	req, err := http.Get(url)
	if err != nil {
		return resp, errors.Wrap(err, "fetch currencies")
	}

	defer req.Body.Close()
	b, err := io.ReadAll(req.Body)
	if err = json.Unmarshal(b, &resp); err != nil {
		return resp, errors.Wrap(err, "read response body")
	}

	return resp, err
}
