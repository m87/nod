package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

func testQueryMixedParameters(t *testing.T, factory RepositoryFactory) {
	repo := createQueryTestRepository(t, factory)

	nodes, err := nod.NewNodeQuery(repo).
		Where(nod.And(
			nod.NodeFields.Kind.Equals("article"),
			nod.NodeFields.Status.Equals("published"),
			nod.KvString("language").Equals("pl"),
			nod.Content("body").Equals("alpha body"),
			nod.Tags().Has("featured"),
		)).
		FindAll()

	require.NoError(t, err)
	requireQueryNodeNames(t, nodes, "alpha")
}
