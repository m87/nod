package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

func testEdgeQueryMixedParameters(t *testing.T, factory RepositoryFactory) {
	repo := createEdgeQueryTestRepository(t, factory)

	edges, err := nod.NewEdgeQuery(repo).
		Where(nod.And(
			nod.EdgeFields.Kind.Equals("dependency"),
			nod.EdgeFields.Status.Equals("active"),
			nod.KvString("language").Equals("pl"),
			nod.Content("body").Equals("alpha body"),
			nod.Tags().Has("featured"),
		)).
		FindAll()

	require.NoError(t, err)
	requireQueryEdgeNames(t, edges, "alpha")
}
