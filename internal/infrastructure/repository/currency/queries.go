package currency

const (
	queryInsertCurrencyPair = `insert into courses (currency_from, currency_to, well) values ($1, $2, $3);`
	queryUpdateCurrencyPair = `update courses set well = $1 where currency_from = $2 and currency_to = $3;`
	queryGetCurrencyPair    = `select currency_from, currency_to, well from courses where currency_from = $1 and currency_to = $2;`
	queryGetAllCurrencyPair = `select currency_from, currency_to from courses;`
)
