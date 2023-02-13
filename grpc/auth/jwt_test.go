package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	token, err := generateToken("user")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	fmt.Println("generated token:", token)
}
