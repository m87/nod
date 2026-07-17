package contract

import "testing"



func testQueries(t *testing.T, factory RepositoryFactory) {
	t.Helper()

	t.Run("FindAllNodes", func(t *testing.T) { testFindAllNodes(t, factory) })

}
