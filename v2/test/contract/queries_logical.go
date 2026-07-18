package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

func testQueryLogicalOperators(t *testing.T, factory RepositoryFactory) {
	repo := createQueryTestRepository(t, factory)

	t.Run("and", func(t *testing.T) {
		nodes, err := nod.NewNodeQuery(repo).
			Where(nod.And(
				nod.CoreFields.Kind.Equals("article"),
				nod.CoreFields.Status.Equals("published"),
			)).
			FindAll()

		require.NoError(t, err)
		requireQueryNodeNames(t, nodes, "alpha")
	})

	t.Run("or", func(t *testing.T) {
		nodes, err := nod.NewNodeQuery(repo).
			Where(nod.Or(
				nod.CoreFields.Kind.Equals("note"),
				nod.CoreFields.Kind.Equals("task"),
			)).
			FindAll()

		require.NoError(t, err)
		requireQueryNodeNames(t, nodes, "gamma", "delta")
	})

	t.Run("and with nested or", func(t *testing.T) {
		nodes, err := nod.NewNodeQuery(repo).
			Where(nod.And(
				nod.Or(
					nod.CoreFields.Kind.Equals("article"),
					nod.CoreFields.Kind.Equals("note"),
				),
				nod.CoreFields.Status.Equals("published"),
				nod.Or(
					nod.Tags().Has("featured"),
					nod.KvString("color").Equals("blue"),
				),
			)).
			FindAll()

		require.NoError(t, err)
		requireQueryNodeNames(t, nodes, "alpha", "gamma")
	})

	t.Run("or with nested and", func(t *testing.T) {
		nodes, err := nod.NewNodeQuery(repo).
			Where(nod.Or(
				nod.And(
					nod.CoreFields.Kind.Equals("article"),
					nod.CoreFields.Status.Equals("draft"),
				),
				nod.And(
					nod.CoreFields.Kind.Equals("note"),
					nod.CoreFields.Status.Equals("published"),
				),
			)).
			FindAll()

		require.NoError(t, err)
		requireQueryNodeNames(t, nodes, "beta", "gamma")
	})
}
