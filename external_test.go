package qp

import (
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MyValuer struct{}

func (v MyValuer) Value() (driver.Value, error) {
	return "NULL", nil
}

func Test_in(t *testing.T) {
	t.Run("No slice args", func(t *testing.T) {
		q, args, err := in("id = ?", "1")
		assert.NoError(t, err)
		assert.Equal(t, "id = ?", q)
		assert.Equal(t, []interface{}{"1"}, args)
	})

	t.Run("Single slice arg", func(t *testing.T) {
		q, args, err := in("id IN (?)", []string{"1", "2"})
		assert.NoError(t, err)
		assert.Equal(t, "id IN (?, ?)", q)
		assert.Equal(t, []interface{}{"1", "2"}, args)
	})

	t.Run("Multiple args with a slice", func(t *testing.T) {
		q, args, err := in("id IN (?) AND status = ?", []string{"1", "2"}, "active")
		assert.NoError(t, err)
		assert.Equal(t, "id IN (?, ?) AND status = ?", q)
		assert.Equal(t, []interface{}{"1", "2", "active"}, args)
	})

	t.Run("Empty slice arg", func(t *testing.T) {
		q, args, err := in("id IN (?)", []string{})
		assert.Error(t, err)
		assert.Equal(t, "", q)
		assert.Nil(t, args)
		assert.EqualError(t, err, "empty slice passed to 'in' query")
	})

	t.Run("Valuer", func(t *testing.T) {
		q, args, err := in("id IN (?)", []sql.NullString{{String: "1", Valid: true}, {String: "2"}})
		assert.NoError(t, err)
		assert.Equal(t, "id IN (?, ?)", q)
		assert.Equal(t, []interface{}{sql.NullString{String: "1", Valid: true}, sql.NullString{String: "2", Valid: false}}, args)
	})

	t.Run("MyValuer", func(t *testing.T) {
		q, args, err := in("id IN (?)", MyValuer{})
		assert.NoError(t, err)
		assert.Equal(t, "id IN (?)", q)
		assert.Equal(t, []interface{}{MyValuer{}}, args)
	})

	t.Run("Too many bindVars", func(t *testing.T) {
		q, args, err := in("id IN (?), id2 = ?", []string{"1", "2"})
		assert.Error(t, err)
		assert.Equal(t, "", q)
		assert.Nil(t, args)
		assert.EqualError(t, err, "number of bindVars exceeds arguments")
	})

	t.Run("Too few bindVars", func(t *testing.T) {
		s := "2"
		sPtr := &s
		q, args, err := in("id = ?", []string{"1", "2"}, sPtr)
		assert.Error(t, err)
		assert.Equal(t, "", q)
		assert.Nil(t, args)
		assert.EqualError(t, err, "number of bindVars less than number arguments")
	})

	t.Run("Skip not slice", func(t *testing.T) {
		q, args, err := in("id IN (?), id2 = ?", "1", []interface{}{"2"})
		assert.NoError(t, err)
		assert.Equal(t, "id IN (?) AND id2 = ?", q)
		assert.Equal(t, []interface{}{"1", "2"}, args)
	})
}
