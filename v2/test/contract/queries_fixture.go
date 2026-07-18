package contract

import (
	"sort"
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

const (
	queryNodeAlphaID = "query-node-alpha"
	queryNodeBetaID  = "query-node-beta"
	queryNodeGammaID = "query-node-gamma"
	queryNodeDeltaID = "query-node-delta"

	queryNamespaceA = "query-namespace-a"
	queryNamespaceB = "query-namespace-b"
)

func createQueryTestRepository(t *testing.T, factory RepositoryFactory) *nod.Repository {
	t.Helper()

	repo := factory(t)
	t.Cleanup(func() {
		require.NoError(t, repo.Close())
	})

	nodes := []*nod.Node{
		{
			Core: nod.NodeCore{
				Id:          queryNodeAlphaID,
				NamespaceId: nod.Ptr(queryNamespaceA),
				Name:        "alpha",
				Kind:        "article",
				Status:      "published",
			},
			KV: map[string]*nod.NodeKV{
				"color":    {Key: "color", ValueText: nod.Ptr("red")},
				"language": {Key: "language", ValueText: nod.Ptr("pl")},
			},
			Content: map[string]*nod.NodeContent{
				"body":    {Key: "body", Value: nod.Ptr("alpha body")},
				"summary": {Key: "summary", Value: nod.Ptr("alpha summary")},
			},
			Tags: []*nod.Tag{
				{Name: "news"},
				{Name: "featured"},
				{Name: "shared"},
			},
		},
		{
			Core: nod.NodeCore{
				Id:          queryNodeBetaID,
				NamespaceId: nod.Ptr(queryNamespaceA),
				ParentId:    nod.Ptr(queryNodeAlphaID),
				Name:        "beta",
				Kind:        "article",
				Status:      "draft",
			},
			KV: map[string]*nod.NodeKV{
				"color":    {Key: "color", ValueText: nod.Ptr("red")},
				"language": {Key: "language", ValueText: nod.Ptr("en")},
				"accent":   {Key: "accent", ValueText: nod.Ptr("blue")},
			},
			Content: map[string]*nod.NodeContent{
				"body":    {Key: "body", Value: nod.Ptr("beta body")},
				"summary": {Key: "summary", Value: nod.Ptr("alpha body")},
			},
			Tags: []*nod.Tag{
				{Name: "tech"},
				{Name: "shared"},
			},
		},
		{
			Core: nod.NodeCore{
				Id:          queryNodeGammaID,
				NamespaceId: nod.Ptr(queryNamespaceB),
				Name:        "gamma",
				Kind:        "note",
				Status:      "published",
			},
			KV: map[string]*nod.NodeKV{
				"color":    {Key: "color", ValueText: nod.Ptr("blue")},
				"language": {Key: "language", ValueText: nod.Ptr("pl")},
			},
			Content: map[string]*nod.NodeContent{
				"body":    {Key: "body", Value: nod.Ptr("gamma body")},
				"summary": {Key: "summary", Value: nod.Ptr("gamma summary")},
			},
			Tags: []*nod.Tag{
				{Name: "news"},
				{Name: "shared"},
			},
		},
		{
			Core: nod.NodeCore{
				Id:          queryNodeDeltaID,
				NamespaceId: nod.Ptr(queryNamespaceB),
				ParentId:    nod.Ptr(queryNodeGammaID),
				Name:        "delta",
				Kind:        "task",
				Status:      "archived",
			},
			KV: map[string]*nod.NodeKV{
				"color":    {Key: "color", ValueText: nod.Ptr("green")},
				"language": {Key: "language", ValueText: nod.Ptr("de")},
			},
			Content: map[string]*nod.NodeContent{
				"body":    {Key: "body", Value: nod.Ptr("delta body")},
				"summary": {Key: "summary", Value: nod.Ptr("delta summary")},
			},
			Tags: []*nod.Tag{
				{Name: "ops"},
			},
		},
	}

	for _, node := range nodes {
		_, err := repo.Nodes().SaveNode(node)
		require.NoError(t, err)
	}

	return repo
}

func requireQueryNodeNames(t *testing.T, nodes []*nod.Node, expected ...string) {
	t.Helper()

	actual := make([]string, 0, len(nodes))
	for _, node := range nodes {
		actual = append(actual, node.Core.Name)
	}

	sort.Strings(actual)
	sort.Strings(expected)
	require.Equal(t, expected, actual)
}
