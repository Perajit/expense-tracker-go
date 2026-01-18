package apperror

import (
	"errors"
)

var (
	ErrDefault                = errors.New("something went wrong")
	ErrInvalidRequest         = errors.New("invalid request")
	ErrNotFound               = errors.New("not found")
	ErrUnauthorized           = errors.New("unauthorized")
	ErrUserDuplication        = errors.New("user already exists")
	ErrInvalidCredentials     = errors.New("invalid username or password")
	ErrInvalidToken           = errors.New("invalid or expired token")
	ErrSecurityContextMissing = errors.New("security context missing")
	ErrRecordDuplication      = errors.New("record already exists")
)
