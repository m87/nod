package contract

import "testing"



func testQueries(t *testing.T, factory RepositoryFactory) {
	t.Helper()

	t.Run("FindAllNodes", func(t *testing.T) { testFindAllNodes(t, factory) })
	t.Run("FindAllNodesWithNoFilter", func(t *testing.T) { testFindAllNodesWithNoFilter(t, factory) })
	t.Run("FindMultipleNodes", func(t *testing.T) { testFindMultipleNodes(t, factory) })
	t.Run("FindByKV", func(t *testing.T) { testFindByKv(t, factory) })
	t.Run("FindByNodeAndKV", func(t *testing.T) { testFindByCoreAndKv(t, factory) })
	t.Run("FullSearch", func(t *testing.T) { testFullSearch(t, factory) })
}
