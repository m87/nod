package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

func testTypedRepositorySaveAndQuery(t *testing.T, factory RepositoryFactory) {
	t.Helper()

	repo := factory(t)
	defer closeRepo(t, repo)

	nod.RegisterMapper(repo.Mappers(), contractModelMapper{})

	typed := nod.As[contractModel](repo)
	model := &contractModel{
		Name: "typed-name",
		Note: "typed-note",
		Tag:  "typed-tag",
	}

	id, err := typed.Save(model)
	require.NoError(t, err)
	require.NotEmpty(t, id)

	found, err := typed.Query().NameEquals("typed-name").KV().Content().Tags().First()
	require.NoError(t, err)
	require.Equal(t, id, found.ID)
	require.Equal(t, "typed-name", found.Name)
	require.Equal(t, "typed-note", found.Note)
	require.Equal(t, "typed-tag", found.Tag)
}
