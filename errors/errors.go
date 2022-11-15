package errors

import "errors"

var (
	// ErrNoRegion No region specified
	ErrNoRegion = errors.New("region parameter is mandatory")

	// ErrInvalidRegion The requested region is not a valid AWS region
	ErrInvalidRegion = errors.New("invalid region requested")
)
