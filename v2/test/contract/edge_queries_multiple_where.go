package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

func testEdgeQueryMultipleWhere(t *testing.T, factory RepositoryFactory) {
	repo := createEdgeQueryTestRepository(t, factory)

	edges, err := nod.NewEdgeQuery(repo).
		Where(nod.EdgeFields.Kind.Equals("dependency")).
		Where(nod.Or(
			nod.EdgeFields.Status.Equals("active"),
			nod.Tags().Has("tech"),
		)).
		Where(nod.KvString("color").Equals("red")).
		FindAll()

	require.NoError(t, err)
	requireQueryEdgeNames(t, edges, "alpha", "beta")
}
