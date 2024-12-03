package service

import "errors"

var (
	ErrOrderAlreadyExists = errors.New("order with this uid already exists")
	ErrOrderNotFound      = errors.New("order not found")
)
