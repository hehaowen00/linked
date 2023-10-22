package items

import "errors"

var (
	ErrMissingUserId = errors.New("missing user id")
	ErrMissingURL    = errors.New("missing url")
	ErrMissingTitle  = errors.New("missing title")
)
