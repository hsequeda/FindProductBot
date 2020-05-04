package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_GetProductsByPattern(t *testing.T) {
	_, err := GetProductsByPattern("carlos3", "ron")
	require.NoError(t, err)
}
