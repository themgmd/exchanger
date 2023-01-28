package currency

import "errors"

var (
	ErrCurrencyPairNotExist     = errors.New("current pair not exists")
	ErrCurrencyPairAlreadyExist = errors.New("current currency pair already exist")
)
