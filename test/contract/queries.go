package contract

import "testing"

func testQueries(t *testing.T, factory RepositoryFactory) {
	t.Helper()

	t.Run("Basic", func(t *testing.T) { testQueryBasic(t, factory) })
	t.Run("CoreFields", func(t *testing.T) { testQueryCoreFields(t, factory) })
	t.Run("Tags", func(t *testing.T) { testQueryTags(t, factory) })
	t.Run("Content", func(t *testing.T) { testQueryContent(t, factory) })
	t.Run("KV", func(t *testing.T) { testQueryKV(t, factory) })
	t.Run("MixedParameters", func(t *testing.T) { testQueryMixedParameters(t, factory) })
	t.Run("LogicalOperators", func(t *testing.T) { testQueryLogicalOperators(t, factory) })
	t.Run("MultipleWhere", func(t *testing.T) { testQueryMultipleWhere(t, factory) })
	t.Run("LazyLoading", func(t *testing.T) { testQueryLazyLoading(t, factory) })
	t.Run("Typed", func(t *testing.T) { testTypedNodeQuery(t, factory) })
}
