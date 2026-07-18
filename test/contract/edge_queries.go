package contract

import "testing"

func testEdgeQueries(t *testing.T, factory RepositoryFactory) {
	t.Helper()

	t.Run("Basic", func(t *testing.T) { testEdgeQueryBasic(t, factory) })
	t.Run("CoreFields", func(t *testing.T) { testEdgeQueryCoreFields(t, factory) })
	t.Run("Tags", func(t *testing.T) { testEdgeQueryTags(t, factory) })
	t.Run("Content", func(t *testing.T) { testEdgeQueryContent(t, factory) })
	t.Run("KV", func(t *testing.T) { testEdgeQueryKV(t, factory) })
	t.Run("MixedParameters", func(t *testing.T) { testEdgeQueryMixedParameters(t, factory) })
	t.Run("LogicalOperators", func(t *testing.T) { testEdgeQueryLogicalOperators(t, factory) })
	t.Run("MultipleWhere", func(t *testing.T) { testEdgeQueryMultipleWhere(t, factory) })
	t.Run("LazyLoading", func(t *testing.T) { testEdgeQueryLazyLoading(t, factory) })
	t.Run("Typed", func(t *testing.T) { testTypedEdgeQuery(t, factory) })
}
