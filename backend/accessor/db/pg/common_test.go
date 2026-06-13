package pg

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
)

func TestIsDuplicateKeyErr(t *testing.T) {
	t.Run("returns constraint name for unique violation", func(t *testing.T) {
		ok, constraint := IsDuplicateKeyErr(&pgconn.PgError{
			Code:           "23505",
			ConstraintName: "note_title_key",
		})

		assert.True(t, ok)
		assert.Equal(t, "note_title_key", constraint)
	})

	t.Run("ignores non-duplicate postgres errors", func(t *testing.T) {
		ok, constraint := IsDuplicateKeyErr(&pgconn.PgError{Code: "23503"})

		assert.False(t, ok)
		assert.Empty(t, constraint)
	})

	t.Run("ignores non-postgres errors", func(t *testing.T) {
		ok, constraint := IsDuplicateKeyErr(errors.New("boom"))

		assert.False(t, ok)
		assert.Empty(t, constraint)
	})
}
