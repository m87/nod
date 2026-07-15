package contract

import (
	"testing"
)

func testNodeCrud(t *testing.T, factory RepositoryFactory) {
	t.Helper()

	t.Run("BasicNodeSave", func(t *testing.T) { testBasicNodeSave(t, factory) })
	t.Run("NodeSaveWithParent", func(t *testing.T) { testNodeSaveWithParent(t, factory) })
	t.Run("NodeDelete", func(t *testing.T) { testNodeDelete(t, factory) })
	t.Run("FullNodeSave", func(t *testing.T) { testFullNodeSave(t, factory) })

}
