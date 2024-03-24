package repo

import (
	"context"
	"exchanger/internal/currency/types"
	"exchanger/pkg/data"
	"exchanger/pkg/pgscan"
	"exchanger/pkg/postgre"
	"fmt"
	"golang.org/x/sync/errgroup"
)

type Currency struct {
	db *postgre.DB
}

func NewCurrency(db *postgre.DB) *Currency {
	return &Currency{
		db: db,
	}
}

func (c Currency) Create(ctx context.Context, pair types.CurrencyPair) error {
	query := `insert into courses (currency_from, currency_to, rate) values (:currency_from, :currency_to, :rate);`
	args := postgre.Args{
		"currency_from": pair.CurrencyFrom,
		"currency_to":   pair.CurrencyTo,
		"rate":          pair.Rate,
	}

	err := pgscan.Exec(ctx, c.db, query, args)
	if err != nil {
		return err
	}

	return nil
}

func (c Currency) CheckExist(ctx context.Context, from, to string) error {
	var exist bool

	query := `
	select
		case
			when count(*) > 0 then true
			else false
		end
	from courses
	where
	    currency_from = :currency_from and
	    currency_to = :currency_to
	;
	`

	err := pgscan.Get(ctx, c.db, &exist, query, postgre.Args{
		"currency_from": from,
		"currency_to":   to,
	})
	if err != nil {
		return err
	}

	if !exist {
		return types.ErrCurrencyPairNotExist
	}

	return nil
}

func (c Currency) Update(ctx context.Context, id int, update data.Json) error {
	sanitizer := data.NewSanitizer()
	update = sanitizer.Sanitize(update, postgre.BaseModelFields...)

	args := postgre.Args(update)

	query := fmt.Sprintf(`update courses set %s where id = :id;`, args.String())
	clear(args)

	args["id"] = id

	err := pgscan.Exec(ctx, c.db, query, args)
	if err != nil {
		return err
	}

	return nil
}

func (c Currency) Get(ctx context.Context, from, to string) (types.CurrencyPair, error) {
	var pair types.CurrencyPair

	query := `select
				id,
				created_at,
				updated_at,
				currency_from,
				currency_to,
				rate
			  from courses
			  where
			    currency_from = :currency_from and
			    currency_to = :currency_to
			  ;
	`

	err := pgscan.Get(ctx, c.db, &pair, query, postgre.Args{
		"currency_from": from,
		"currency_to":   to,
	})
	if err != nil {
		return types.CurrencyPair{}, err
	}

	return pair, nil
}

func (c Currency) GetById(ctx context.Context, id int) (types.CurrencyPair, error) {
	var pair types.CurrencyPair

	query := `select
				id,
				created_at,
				updated_at,
				currency_from,
				currency_to,
				rate
			  from courses
			  where
			    id = :id
			  ;
	`

	err := pgscan.Get(ctx, c.db, &pair, query, postgre.Args{
		"id": id,
	})
	if err != nil {
		return types.CurrencyPair{}, err
	}

	return pair, nil
}

func (c Currency) GetRate(ctx context.Context, from, to string) (float64, error) {
	var rate float64

	query := `select
				rate
			  from courses
			  where
			    currency_from = :currency_from and
			    currency_to = :currency_to
			  ;
	`

	err := pgscan.Get(ctx, c.db, &rate, query, postgre.Args{
		"currency_from": from,
		"currency_to":   to,
	})
	if err != nil {
		return 0, err
	}

	return rate, nil
}

func (c Currency) List(ctx context.Context, limit, offset int) ([]types.CurrencyPair, int, error) {
	var (
		pairs = make([]types.CurrencyPair, 0)
		total int
		eg, _ = errgroup.WithContext(ctx)
	)

	countQuery := `select count(*) from courses;`
	selectQuery := `select
				id,
				created_at,
				updated_at,
				currency_from,
				currency_to,
				rate
			  from courses
			  limit :limit
			  offset :offset;
	`

	eg.Go(func() error {
		return pgscan.Get(ctx, c.db, &total, countQuery, postgre.Args{})
	})

	eg.Go(func() error {
		return pgscan.Select(ctx, c.db, &pairs, selectQuery, postgre.Args{
			"limit":  limit,
			"offset": offset,
		})
	})

	err := eg.Wait()
	if err != nil {
		return pairs, 0, err
	}

	return pairs, total, nil
}
