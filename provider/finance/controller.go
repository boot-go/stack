package finance

import "github.com/piquette/finance-go"

type Controller interface {
	Quote(symbol string) (*finance.Quote, error)
}
