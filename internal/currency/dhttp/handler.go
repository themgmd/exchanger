package dhttp

import (
	"context"
	"exchanger/internal/currency/types"
	"exchanger/pkg/errors"
	httplib "exchanger/pkg/http"
	"exchanger/pkg/pagination"
	"net/http"
	"strconv"
)

type Service interface {
	CreatePair(ctx context.Context, from, to string) (types.CurrencyPair, error)
	Exchange(ctx context.Context, from, to string, amount float64) (float64, error)
	UpdateRate(ctx context.Context, id int, rate float64) error
	GetRate(ctx context.Context, from, to string) (types.CurrencyPair, error)
	List(ctx context.Context, pag pagination.Pagination) ([]types.CurrencyPair, int, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h Handler) CreatePair(w http.ResponseWriter, r *http.Request) {
	var req types.CurrencyPairsRequest

	err := httplib.ReadBody(r.Body, &req)
	if err != nil {
		httplib.NewBadRequestResponse(w, err)
		return
	}

	pair, err := h.service.CreatePair(r.Context(), req.CurrencyFrom, req.CurrencyTo)
	if err != nil {
		httplib.NewInternalServerErrorResponse(w, err)
		return
	}

	httplib.NewSuccessResponse(w, pair)
}

func (h Handler) Exchange(w http.ResponseWriter, r *http.Request) {
	var req types.ExchangeRequest

	err := httplib.ReadBody(r.Body, &req)
	if err != nil {
		httplib.NewBadRequestResponse(w, err)
		return
	}

	amount, err := h.service.Exchange(r.Context(), req.CurrencyFrom, req.CurrencyTo, req.Amount)
	if err != nil {
		httplib.NewInternalServerErrorResponse(w, err)
		return
	}

	httplib.NewSuccessResponse(w, amount)
}

func (h Handler) GetRate(w http.ResponseWriter, r *http.Request) {
	var (
		currencyFrom = r.PathValue("currency")
		currencyTo   = r.URL.Query().Get("currencyTo")
	)

	if currencyTo == "" {
		httplib.NewBadRequestResponse(w, errors.New("currencyTo could not be empty"))
		return
	}

	amount, err := h.service.GetRate(r.Context(), currencyFrom, currencyTo)
	if err != nil {
		httplib.NewInternalServerErrorResponse(w, err)
		return
	}

	httplib.NewSuccessResponse(w, amount)
}

func (h Handler) List(w http.ResponseWriter, r *http.Request) {
	var (
		err         error
		queryParams = r.URL.Query()
		req         pagination.RequestPagination
	)

	req.CurrentPage, err = strconv.Atoi(queryParams.Get("currentPage"))
	if err != nil {
		httplib.NewInternalServerErrorResponse(w, err)
		return
	}

	req.PerPage, err = strconv.Atoi(queryParams.Get("perPage"))
	if err != nil {
		httplib.NewInternalServerErrorResponse(w, err)
		return
	}

	pairs, total, err := h.service.List(r.Context(), *req.ToPagination())
	if err != nil {
		httplib.NewInternalServerErrorResponse(w, err)
		return
	}

	httplib.NewListResponse(w,
		*pagination.NewResponsePagination(req, total),
		pairs,
	)
}
