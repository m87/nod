package contract

import "testing"

func testNodeTyped(t *testing.T, factory RepositoryFactory) {
	t.Helper()

	t.Run("CodecSave", func(t *testing.T) { testCodecSave(t, factory) })

}
