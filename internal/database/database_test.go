package database

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("Should: not return error", func(t *testing.T) {
		_, s := New(func() (*gorm.DB, error){return nil, errorDataBase})
		assert.Error(t, s)
	})
}