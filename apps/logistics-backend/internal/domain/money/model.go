package money

import (
	"errors"
	"fmt"
	"strings"
)

type Money struct {
	Amount   int64  `db:"amount" json:"amount"`
	Currency string `db:"currency" json:"currency"`
}

var currencySymbols = map[string]string{
	"USD": "$",
	"EUR": "€",
	"GBP": "£",
	"KES": "KSh",
}

// New creates a Money value from a float and currency
func New(amount float64, currency string) Money {
	return Money{
		Amount:   int64(amount * 100),
		Currency: strings.ToUpper(currency),
	}
}

// FromCents creates a Money value from cents
func FromCents(cents int64, currency string) Money {
	return Money{
		Amount:   cents,
		Currency: strings.ToUpper(currency),
	}
}

// String formats as "KSh 1,234.56" or "$1,234.56"
func (m Money) String() string {
	symbol, exists := currencySymbols[m.Currency]
	if !exists {
		symbol = m.Currency // fallback to currency code
	}

	return fmt.Sprintf("%s %.2f", symbol, float64(m.Amount)/100)
}

// Add returns a new Money with summed amounts
func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, errors.New("currency mismatch")
	}
	return Money{
		Amount:   m.Amount + other.Amount,
		Currency: m.Currency,
	}, nil
}

// Multiply returns a new Money multiplied by a scalar
func (m Money) Multiply(factor int64) Money {
	return Money{
		Amount:   m.Amount * factor,
		Currency: m.Currency,
	}
}
