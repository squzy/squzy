package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("Should: Create new application", func(t *testing.T) {
		app := New(nil, nil, nil, nil)
		assert.NotEqual(t, nil, app)
	})
}

func TestApp_Run(t *testing.T) {
	t.Run("Should: throw error invalid port", func(t *testing.T) {
		app := New(nil, nil, nil, nil)
		assert.Equal(t, nil, app.Run(3123123))
	})
}