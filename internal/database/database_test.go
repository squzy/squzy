package database

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("Should: not return nil", func(t *testing.T) {
		s, err := New(nil)
		assert.NotNil(t, err)
		assert.Nil(t, s)
	})
}

func TestNewPostgres(t *testing.T) {
	t.Run("Should: not return nil", func(t *testing.T) {
		db, _ := gorm.Open("postgres", "user=gorm password=gorm DB.name=gorm port=9920 sslmode=disable")
		err := os.Setenv("DB_TYPE", "postgres")
		if err != nil {
			assert.Fail(t, err.Error())
		}

		s, err := New(db)
		assert.Nil(t, err)
		assert.NotNil(t, s)
	})
}

func TestNewPostgresErr(t *testing.T) {
	t.Run("Should: not return nil", func(t *testing.T) {
		db, _, err := sqlmock.New()
		if err != nil {
			assert.Fail(t, err.Error())
		}
		err = os.Setenv("DB_TYPE", "postgres")
		if err != nil {
			assert.Fail(t, err.Error())
		}
		s, err := New(db)
		assert.NotNil(t, err)
		assert.Nil(t, s)
	})
}

func TestNewClichouse(t *testing.T) {
	t.Run("Should: not return nil", func(t *testing.T) {
		db, _, err := sqlmock.New()
		if err != nil {
			assert.Fail(t, err.Error())
		}
		err = os.Setenv("DB_TYPE", "clickhouse")
		if err != nil {
			assert.Fail(t, err.Error())
		}

		s, err := New(db)
		assert.Nil(t, err)
		assert.NotNil(t, s)
	})
}

func TestNewClichouseErr(t *testing.T) {
	t.Run("Should: not return nil", func(t *testing.T) {
		db, _ := gorm.Open("postgres", "user=gorm password=gorm DB.name=gorm port=9920 sslmode=disable")

		err := os.Setenv("DB_TYPE", "clickhouse")
		if err != nil {
			assert.Fail(t, err.Error())
		}

		s, err := New(db)
		assert.NotNil(t, err)
		assert.Nil(t, s)
	})
}
