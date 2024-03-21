package dhttp

import "net/http"

type Currency struct {
	handler *Handler
}

func NewCurrency(service Service) *Currency {
	return &Currency{
		handler: NewHandler(service),
	}
}

func (c Currency) SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /currency", nil)
	mux.HandleFunc("GET /currency", nil)
}
