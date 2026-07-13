package contract

import (
	"testing"

	"github.com/m87/nod"
)

type RepositoryFactory func(t *testing.T) *nod.Repository

func RunRepositoryContractTests(t *testing.T, factory RepositoryFactory) {
	t.Helper()

	t.Run("NodeCrud", func(t *testing.T) { testNodeCrud(t, factory) })
	t.Run("EdgeCrud", func(t *testing.T) { testEdgeCrud(t, factory) })
}
