package errors

import "errors"

var (

	// ErrInvalidRegion The requested region is not a valid AWS region
	ErrInvalidRegion = errors.New("invalid region requested")
)
