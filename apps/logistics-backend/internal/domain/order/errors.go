package order

import "errors"

var (
	ErrorOutOfStock           = errors.New("product out of stock")
	ErrorInvalidQuantity      = errors.New("invalid quantity")
	ErrorQuantityExceedsStock = errors.New("ordered quantity exceeds available stock")
)
