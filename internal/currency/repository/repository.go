package repository

import (
	"context"
	"database/sql"
	"errors"
	"exchanger/internal/currency"
	"exchanger/internal/models"
	"exchanger/pkg/database/postgres"
	"fmt"
	sq "github.com/Masterminds/squirrel"
)

const (
	scheme      = "currencies"
	table       = "courses"
	tableScheme = scheme + "." + table
)

type repository struct {
	client       postgres.Client
	queryBuilder sq.StatementBuilderType
}

func New(client postgres.Client) currency.Repository {
	return &repository{
		queryBuilder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		client:       client,
	}
}

func (r repository) Insert(ctx context.Context, pair models.CurrencyPair) error {
	query, args, err := r.queryBuilder.
		Insert(tableScheme).
		Columns(columnCurrencyFrom, columnCurrencyTo, columnRate).
		Values(pair.CurrencyFrom, pair.CurrencyTo, pair.Rate).
		ToSql()

	if err != nil {
		return fmt.Errorf("queryBuilder: %w", err)
	}

	exec, err := r.client.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf(" r.client.Exec: %w", err)
	}

	if !exec.Insert() {
		return errors.New("insert in db failed")
	}

	return nil
}

func (r repository) Update(ctx context.Context, pair models.CurrencyPair) error {
	query, args, err := r.queryBuilder.
		Insert(tableScheme).
		Columns(columnCurrencyFrom, columnCurrencyTo, columnRate).
		Values(pair.CurrencyFrom, pair.CurrencyTo, pair.Rate).
		ToSql()

	if err != nil {
		return fmt.Errorf("queryBuilder: %w", err)
	}

	exec, err := r.client.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf(" r.client.Exec: %w", err)
	}

	if !exec.Update() || exec.RowsAffected() == 0 {
		return errors.New("update in db failed")
	}

	return nil
}

func (r repository) CheckExist(ctx context.Context, params models.CurrencyParams) (bool, error) {
	var id int
	query, args, err := r.queryBuilder.
		Select(columnId).
		From(tableScheme).
		Where(sq.Eq{
			columnCurrencyFrom: params.CurrencyFrom,
			columnCurrencyTo:   params.CurrencyTo,
		}).
		ToSql()
	if err != nil {
		return false, fmt.Errorf("queryBuilder: %w", err)
	}

	row := r.client.QueryRow(ctx, query, args...)
	if err := row.Scan(&id); err != nil && errors.Is(err, sql.ErrNoRows) {
		return false, fmt.Errorf(" r.client.Exec: %w", err)
	}

	if id != 0 {
		return true, nil
	}

	return false, currency.ErrCurrencyPairNotExist
}

func (r repository) GetRate(ctx context.Context, params models.CurrencyParams) (float64, error) {
	var rate float64
	query, args, err := r.queryBuilder.
		Select(columnRate).
		From(tableScheme).
		Where(sq.Eq{
			columnCurrencyFrom: params.CurrencyFrom,
			columnCurrencyTo:   params.CurrencyTo,
		}).
		ToSql()
	if err != nil {
		return 0, fmt.Errorf("queryBuilder: %w", err)
	}

	row := r.client.QueryRow(ctx, query, args...)
	if err := row.Scan(&rate); err != nil {
		return 0, fmt.Errorf(" r.client.Exec: %w", err)
	}

	return rate, nil
}

func (r repository) Get(ctx context.Context, params models.CurrencyParams) (*models.CurrencyPair, error) {
	var model models.CurrencyPair
	query, args, err := r.queryBuilder.
		Select(columnCurrencyFrom, columnCurrencyTo, columnRate).
		From(tableScheme).
		Where(sq.Eq{
			columnCurrencyFrom: params.CurrencyFrom,
			columnCurrencyTo:   params.CurrencyTo,
		}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("queryBuilder: %w", err)
	}

	row := r.client.QueryRow(ctx, query, args...)
	if err := row.Scan(&model.CurrencyFrom, &model.CurrencyTo, &model.Rate); err != nil {
		return nil, fmt.Errorf("r.client.QueryRow: %w", err)
	}

	return &model, nil
}

func (r repository) Select(ctx context.Context, limit, offset int) ([]models.CurrencyPair, error) {
	var result []models.CurrencyPair
	sqlQuery := r.queryBuilder.
		Select(columnCurrencyFrom, columnCurrencyTo, columnRate).
		From(tableScheme)

	if limit != 0 {
		sqlQuery = sqlQuery.Limit(uint64(limit))
	}

	if offset != 0 {
		sqlQuery = sqlQuery.Offset(uint64(offset))
	}

	query, args, err := sqlQuery.ToSql()
	if err != nil {
		return result, fmt.Errorf("queryBuilder: %w", err)
	}

	rows, err := r.client.Query(ctx, query, args...)
	defer rows.Close()
	if err != nil {
		return result, fmt.Errorf(" r.client.Query: %w", err)
	}

	for rows.Next() {
		var res models.CurrencyPair
		if err = rows.Scan(&res.CurrencyFrom, &res.CurrencyTo, &res.Rate); err != nil {
			return result, fmt.Errorf("rows.Scan: %w", err)
		}
		result = append(result, res)
	}

	if err = rows.Err(); err != nil {
		return result, fmt.Errorf("rows.Err: %w", err)
	}

	return result, nil
}
