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

	typed := nod.NewTypedRepository[contractModel](repo)
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

type contractModelView struct {
	ID   string
	Name string
}

type contractModelViewMapper struct{}

func (contractModelViewMapper) ToNode(model *contractModelView) (*nod.Node, error) {
	return &nod.Node{Core: nod.NodeCore{
		Id: model.ID, Name: model.Name, Kind: "contract-kind", Status: "active",
	}}, nil
}

func (contractModelViewMapper) FromNode(node *nod.Node) (*contractModelView, error) {
	return &contractModelView{ID: node.Core.Id, Name: node.Core.Name}, nil
}

func (contractModelViewMapper) IsApplicable(node *nod.Node) bool {
	return node.Core.Kind == "contract-kind"
}

func testTypedParentsAndConversion(t *testing.T, factory RepositoryFactory) {
	t.Helper()

	repo := factory(t)
	defer closeRepo(t, repo)

	nod.RegisterMapper(repo.Mappers(), contractModelMapper{})
	nod.RegisterMapper(repo.Mappers(), contractModelViewMapper{})

	parentID, err := repo.Save(&nod.Node{Core: nod.NodeCore{
		Id: "typed-parent", Name: "matching-parent", Kind: "contract-kind", Status: "active",
	}})
	require.NoError(t, err)

	for _, child := range []*nod.Node{
		{Core: nod.NodeCore{Id: "typed-child-1", ParentId: ptr(parentID), Name: "wanted-child-1", Kind: "contract-kind", Status: "active"}},
		{Core: nod.NodeCore{Id: "typed-child-2", ParentId: ptr(parentID), Name: "wanted-child-2", Kind: "contract-kind", Status: "active"}},
		{Core: nod.NodeCore{Id: "typed-child-3", Name: "wanted-root", Kind: "contract-kind", Status: "active"}},
	} {
		_, err = repo.Save(child)
		require.NoError(t, err)
	}

	models := nod.NewTypedRepository[contractModel](repo)
	views := nod.NewTypedRepository[contractModelView](repo)
	children := models.Query().NameStartsWith("wanted-child")

	sameType, err := children.Parents().NameEquals("matching-parent").List()
	require.NoError(t, err)
	require.Len(t, sameType, 1)
	require.Equal(t, parentID, sameType[0].ID)

	otherType, err := nod.QueryAs[contractModelView](children.Parents()).List()
	require.NoError(t, err)
	require.Equal(t, []*contractModelView{{ID: parentID, Name: "matching-parent"}}, otherType)

	rawParent, err := repo.Query().NodeId(parentID).First()
	require.NoError(t, err)

	full, err := models.NodeAs(rawParent)
	require.NoError(t, err)
	require.Equal(t, parentID, full.ID)

	view, err := views.NodeAs(rawParent)
	require.NoError(t, err)
	require.Equal(t, &contractModelView{ID: parentID, Name: "matching-parent"}, view)

	viewFromFunction, err := nod.NodeAs[contractModelView](repo, rawParent)
	require.NoError(t, err)
	require.Equal(t, view, viewFromFunction)

	notApplicable := &nod.Node{Core: nod.NodeCore{Id: "other", Kind: "other-kind"}}
	_, err = views.NodeAs(notApplicable)
	require.ErrorContains(t, err, "not applicable")
}
