package stderrs

import (
	"errors"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrTenderNotFound  = errors.New("tender not found")
	ErrNotEnoughRights = errors.New("not enough rights")
	ErrUpdValue        = errors.New("not enough value to update")
)
