package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigConnStr(t *testing.T) {
	t.Parallel()

	cfg := &Config{
		DbUser:     "user",
		DbPassword: "pwd",
		DbHost:     "localhost",
		DbPort:     "5432",
		DbName:     "app",
	}

	assert.Equal(t, "postgres://user:pwd@localhost:5432/app?sslmode=disable", cfg.ConnStr())
}
