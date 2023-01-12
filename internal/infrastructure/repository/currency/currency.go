package currency

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/onemgvv/exchanger/internal/domain/entity"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db}
}

func (r Repository) Insert(pair entity.CurrencyPair) error {
	row := r.db.QueryRow(queryInsertCurrencyPair, pair.CurrencyFrom, pair.CurrencyTo, pair.Well)
	if !errors.Is(row.Err(), sql.ErrNoRows) {
		return row.Err()
	}

	return nil
}

func (r Repository) Update(pair entity.CurrencyPair) error {
	row := r.db.QueryRow(queryUpdateCurrencyPair, pair.Well, pair.CurrencyFrom, pair.CurrencyTo)
	if !errors.Is(row.Err(), sql.ErrNoRows) {
		return row.Err()
	}

	return nil
}

func (r Repository) Get(find entity.CurrencyPairParams) *entity.CurrencyPair {
	var result entity.CurrencyPair
	err := r.db.Get(&result, queryGetCurrencyPair, find.CurrencyFrom, find.CurrencyTo)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}

	return &result
}

func (r Repository) Select() []entity.CurrencyPairParams {
	result := make([]entity.CurrencyPairParams, 0)
	if err := r.db.Get(&result, queryGetAllCurrencyPair); err != nil {
		return result
	}

	return result
}
