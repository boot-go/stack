package finance

import (
	"github.com/boot-go/boot"
	"github.com/piquette/finance-go"
	"github.com/piquette/finance-go/quote"
)

type component struct {
}

func (c *component) Init() error {
	return nil
}

func (c *component) Quote(symbol string) (*finance.Quote, error) {
	quoteResult, err := quote.Get(symbol)
	if err != nil {
		return nil, err
	}
	return quoteResult, nil
}

func init() {
	boot.Register(func() boot.Component {
		return &component{}
	})
}
