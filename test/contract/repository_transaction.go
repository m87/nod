package contract

import (
	"errors"
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func testRepositoryTransaction(t *testing.T, factory RepositoryFactory) {
	repo := factory(t)
	t.Cleanup(func() {
		require.NoError(t, repo.Close())
	})

	err := repo.Transaction(func(txRepository *nod.Repository) error {
		require.Same(t, repo.Adapters(), txRepository.Adapters())
		_, err := txRepository.Nodes().SaveNode(&nod.Node{
			Core: nod.NodeCore{Id: "transaction-commit", Name: "committed", Kind: "test"},
		})
		return err
	})
	require.NoError(t, err)

	committed, err := repo.Nodes().GetNode("transaction-commit")
	require.NoError(t, err)
	require.Equal(t, "committed", committed.Core.Name)

	wantErr := errors.New("rollback transaction")
	err = repo.Transaction(func(txRepository *nod.Repository) error {
		_, err := txRepository.Nodes().SaveNode(&nod.Node{
			Core: nod.NodeCore{Id: "transaction-rollback", Name: "rolled back", Kind: "test"},
		})
		require.NoError(t, err)
		return wantErr
	})
	require.ErrorIs(t, err, wantErr)

	_, err = repo.Nodes().GetNode("transaction-rollback")
	require.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
