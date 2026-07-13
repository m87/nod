package contract

import (
	"testing"

	"github.com/m87/nod"
)

type RepositoryFactory func(t *testing.T) *nod.Repository

// Runs the full repository contract suite.
func RunRepositoryContractTests(t *testing.T, factory RepositoryFactory) {
	t.Helper()

	t.Run("SaveAndQueryFullModel", func(t *testing.T) { testSaveAndQueryFullModel(t, factory) })
	t.Run("EdgeAndEdgeKV", func(t *testing.T) { testEdgeAndEdgeKV(t, factory) })
	t.Run("Constraints", func(t *testing.T) { testConstraints(t, factory) })
	t.Run("Migration", func(t *testing.T) { testMigration(t, factory) })
	t.Run("TypedRepositorySaveAndQuery", func(t *testing.T) { testTypedRepositorySaveAndQuery(t, factory) })
	t.Run("TypedParentsAndConversion", func(t *testing.T) { testTypedParentsAndConversion(t, factory) })
	t.Run("RepositoryClose", func(t *testing.T) { testRepositoryClose(t, factory) })
}
