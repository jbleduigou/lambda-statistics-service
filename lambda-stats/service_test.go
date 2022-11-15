package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLambdaFunctionsShouldReturnPouet(t *testing.T) {
	s, err := NewLambdaService("us-east-1")

	assert.Nil(t, err)

	stats, err := s.GetLambdaFunctions(context.Background())

	assert.Nil(t, err)

	assert.Len(t, stats, 1)
	assert.Equal(t, "pouet", stats[0])
}
