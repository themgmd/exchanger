package types

import "exchanger/pkg/errors"

var (
	ErrCurrencyPairNotExist     = errors.New("current pair not exists")
	ErrCurrencyPairAlreadyExist = errors.New("current currency_d pair already exist")
)
