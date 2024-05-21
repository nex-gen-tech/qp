package qp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_stringInSlice(t *testing.T) {
	t.Run("RETURN FALSE", func(t *testing.T) {
		assert.Equal(t, false, stringInSlice("", nil))
	})
}
