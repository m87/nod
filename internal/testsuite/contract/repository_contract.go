package contract

import (
	"testing"

	"github.com/m87/nod"
)

type RepositoryFactory func(t *testing.T) *nod.Repository

// RunRepositoryContractTests runs the full repository contract suite. The
// suite is implemented as a set of smaller test functions split across
// files to keep the code easy to read while preserving a single entrypoint
// for adapter tests.
func RunRepositoryContractTests(t *testing.T, factory RepositoryFactory) {
	t.Helper()

	t.Run("SaveAndQueryFullModel", func(t *testing.T) { testSaveAndQueryFullModel(t, factory) })
	t.Run("Constraints", func(t *testing.T) { testConstraints(t, factory) })
	t.Run("Migration", func(t *testing.T) { testMigration(t, factory) })
	t.Run("TagRepositoryDelete", func(t *testing.T) { testTagRepositoryDelete(t, factory) })
	t.Run("TypedRepositorySaveAndQuery", func(t *testing.T) { testTypedRepositorySaveAndQuery(t, factory) })
	t.Run("RepositoryClose", func(t *testing.T) { testRepositoryClose(t, factory) })
}
