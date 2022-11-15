package errors

import "errors"

var (
	// ErrNoRegion No region specified
	ErrNoRegion = errors.New("region parameter is mandatory")

	// ErrInvalidRegion The requested region is not a valid AWS region
	ErrInvalidRegion = errors.New("invalid region requested")

	// ErrNoRuntime No runtime specified
	ErrNoRuntime = errors.New("runtime parameter is mandatory")

	// ErrInvalidRuntime The requested runtime is not a valid AWS runtime
	ErrInvalidRuntime = errors.New("invalid runtime requested")
)
