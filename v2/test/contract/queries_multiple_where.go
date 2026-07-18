package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

func testQueryMultipleWhere(t *testing.T, factory RepositoryFactory) {
	repo := createQueryTestRepository(t, factory)

	nodes, err := nod.NewNodeQuery(repo).
		Where(nod.NodeFields.Kind.Equals("article")).
		Where(nod.Or(
			nod.NodeFields.Status.Equals("published"),
			nod.Tags().Has("tech"),
		)).
		Where(nod.KvString("color").Equals("red")).
		FindAll()

	require.NoError(t, err)
	requireQueryNodeNames(t, nodes, "alpha", "beta")
}
