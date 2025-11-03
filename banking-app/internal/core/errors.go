package core

import "errors"

var (
	ErrNotFound          = errors.New("resource not found")
	ErrInvalidInput      = errors.New("invalid input")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrAccountClosed     = errors.New("account is closed")
	ErrAccountSuspended  = errors.New("account is suspended")
	ErrDuplicateEntry    = errors.New("duplicate entry")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrInternal          = errors.New("internal server error")
)