package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

func testEdgeQueryLogicalOperators(t *testing.T, factory RepositoryFactory) {
	repo := createEdgeQueryTestRepository(t, factory)

	t.Run("and", func(t *testing.T) {
		edges, err := nod.NewEdgeQuery(repo).
			Where(nod.And(
				nod.EdgeFields.Kind.Equals("dependency"),
				nod.EdgeFields.Status.Equals("active"),
			)).
			FindAll()

		require.NoError(t, err)
		requireQueryEdgeNames(t, edges, "alpha")
	})

	t.Run("or", func(t *testing.T) {
		edges, err := nod.NewEdgeQuery(repo).
			Where(nod.Or(
				nod.EdgeFields.Kind.Equals("reference"),
				nod.EdgeFields.Kind.Equals("ownership"),
			)).
			FindAll()

		require.NoError(t, err)
		requireQueryEdgeNames(t, edges, "gamma", "delta")
	})

	t.Run("and with nested or", func(t *testing.T) {
		edges, err := nod.NewEdgeQuery(repo).
			Where(nod.And(
				nod.Or(
					nod.EdgeFields.Kind.Equals("dependency"),
					nod.EdgeFields.Kind.Equals("reference"),
				),
				nod.EdgeFields.Status.Equals("active"),
				nod.Or(
					nod.Tags().Has("featured"),
					nod.KvString("color").Equals("blue"),
				),
			)).
			FindAll()

		require.NoError(t, err)
		requireQueryEdgeNames(t, edges, "alpha", "gamma")
	})

	t.Run("or with nested and", func(t *testing.T) {
		edges, err := nod.NewEdgeQuery(repo).
			Where(nod.Or(
				nod.And(
					nod.EdgeFields.Kind.Equals("dependency"),
					nod.EdgeFields.Status.Equals("inactive"),
				),
				nod.And(
					nod.EdgeFields.Kind.Equals("reference"),
					nod.EdgeFields.Status.Equals("active"),
				),
			)).
			FindAll()

		require.NoError(t, err)
		requireQueryEdgeNames(t, edges, "beta", "gamma")
	})
}
