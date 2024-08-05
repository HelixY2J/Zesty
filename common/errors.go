package common

import "errors"

var (
	ErrNoItems = errors.New("items should have atleast one item")
)
