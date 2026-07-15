package contract

import "testing"

func testEdgeTyped(t *testing.T, factory RepositoryFactory) {
	t.Helper()

	t.Run("CustomEdgeCodec", func(t *testing.T) { testCustomEdgeCodec(t, factory) })
}
