package postgres_test

import (
	"testing"

	"github.com/gosom/goappbuild"
	"github.com/gosom/goappbuild/postgres"
	"github.com/stretchr/testify/require"
)

func Test_postgresQ(t *testing.T) {
	t.Run("test equal", func(t *testing.T) {
		q := goappbuild.Q{}.
			Schema("test").
			Table("users").
			Select("id", "name", "email").
			Equal("id", 1).
			NotEqual("department", "engineering").
			GreaterThanOrEqual("age", 18).
			LessThanOrEqual("age", 65).
			Null("deleted_at").
			NotNull("activated_at").
			StartsWith("name", "John").
			EndsWith("name", "Smith")

		builder := postgres.NewPostgresQ(q)

		sql, args, err := builder.Build()
		require.NoError(t, err)

		expected := `SELECT "id", "name", "email" FROM "test"."users" WHERE "id" = $1 AND "department" != $2 AND "age" >= $3 AND "age" <= $4 AND "deleted_at" IS NULL  AND "activated_at" IS NOT NULL  AND "name" LIKE $5 AND "name" LIKE $6`

		require.Equal(t, expected, sql)
		require.Equal(t, 6, len(args))
		require.Equal(t, 1, args[0])
		require.Equal(t, "engineering", args[1])
		require.Equal(t, 18, args[2])
		require.Equal(t, 65, args[3])
		require.Equal(t, "John%", args[4])
		require.Equal(t, "%Smith", args[5])
	})
}
