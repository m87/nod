package contract

import (
	"sort"
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

const (
	queryEdgeAlphaID = "query-edge-alpha"
	queryEdgeBetaID  = "query-edge-beta"
	queryEdgeGammaID = "query-edge-gamma"
	queryEdgeDeltaID = "query-edge-delta"

	queryEdgeSourceAID = "query-edge-source-a"
	queryEdgeSourceBID = "query-edge-source-b"
	queryEdgeTargetAID = "query-edge-target-a"
	queryEdgeTargetBID = "query-edge-target-b"
)

func createEdgeQueryTestRepository(t *testing.T, factory RepositoryFactory) *nod.Repository {
	t.Helper()

	repo := factory(t)
	t.Cleanup(func() {
		require.NoError(t, repo.Close())
	})

	for _, node := range []*nod.Node{
		{Core: nod.NodeCore{Id: queryEdgeSourceAID, Name: "source-a", Kind: "endpoint"}},
		{Core: nod.NodeCore{Id: queryEdgeSourceBID, Name: "source-b", Kind: "endpoint"}},
		{Core: nod.NodeCore{Id: queryEdgeTargetAID, Name: "target-a", Kind: "endpoint"}},
		{Core: nod.NodeCore{Id: queryEdgeTargetBID, Name: "target-b", Kind: "endpoint"}},
	} {
		_, err := repo.Nodes().SaveNode(node)
		require.NoError(t, err)
	}

	edges := []*nod.Edge{
		{
			Core: nod.EdgeCore{
				Id:          queryEdgeAlphaID,
				NamespaceId: nod.Ptr(queryNamespaceA),
				SourceId:    queryEdgeSourceAID,
				TargetId:    queryEdgeTargetAID,
				Name:        "alpha",
				Kind:        "dependency",
				Status:      "active",
			},
			KV: map[string]*nod.EdgeKV{
				"color":    {Key: "color", ValueText: nod.Ptr("red")},
				"language": {Key: "language", ValueText: nod.Ptr("pl")},
			},
			Content: map[string]*nod.EdgeContent{
				"body":    {Key: "body", Value: nod.Ptr("alpha body")},
				"summary": {Key: "summary", Value: nod.Ptr("alpha summary")},
			},
			Tags: []*nod.Tag{{Name: "news"}, {Name: "featured"}, {Name: "shared"}},
		},
		{
			Core: nod.EdgeCore{
				Id:          queryEdgeBetaID,
				NamespaceId: nod.Ptr(queryNamespaceA),
				SourceId:    queryEdgeSourceAID,
				TargetId:    queryEdgeTargetBID,
				Name:        "beta",
				Kind:        "dependency",
				Status:      "inactive",
			},
			KV: map[string]*nod.EdgeKV{
				"color":    {Key: "color", ValueText: nod.Ptr("red")},
				"language": {Key: "language", ValueText: nod.Ptr("en")},
				"accent":   {Key: "accent", ValueText: nod.Ptr("blue")},
			},
			Content: map[string]*nod.EdgeContent{
				"body":    {Key: "body", Value: nod.Ptr("beta body")},
				"summary": {Key: "summary", Value: nod.Ptr("alpha body")},
			},
			Tags: []*nod.Tag{{Name: "tech"}, {Name: "shared"}},
		},
		{
			Core: nod.EdgeCore{
				Id:          queryEdgeGammaID,
				NamespaceId: nod.Ptr(queryNamespaceB),
				SourceId:    queryEdgeSourceBID,
				TargetId:    queryEdgeTargetAID,
				Name:        "gamma",
				Kind:        "reference",
				Status:      "active",
			},
			KV: map[string]*nod.EdgeKV{
				"color":    {Key: "color", ValueText: nod.Ptr("blue")},
				"language": {Key: "language", ValueText: nod.Ptr("pl")},
			},
			Content: map[string]*nod.EdgeContent{
				"body":    {Key: "body", Value: nod.Ptr("gamma body")},
				"summary": {Key: "summary", Value: nod.Ptr("gamma summary")},
			},
			Tags: []*nod.Tag{{Name: "news"}, {Name: "shared"}},
		},
		{
			Core: nod.EdgeCore{
				Id:          queryEdgeDeltaID,
				NamespaceId: nod.Ptr(queryNamespaceB),
				SourceId:    queryEdgeSourceBID,
				TargetId:    queryEdgeTargetBID,
				Name:        "delta",
				Kind:        "ownership",
				Status:      "archived",
			},
			KV: map[string]*nod.EdgeKV{
				"color":    {Key: "color", ValueText: nod.Ptr("green")},
				"language": {Key: "language", ValueText: nod.Ptr("de")},
			},
			Content: map[string]*nod.EdgeContent{
				"body":    {Key: "body", Value: nod.Ptr("delta body")},
				"summary": {Key: "summary", Value: nod.Ptr("delta summary")},
			},
			Tags: []*nod.Tag{{Name: "ops"}},
		},
	}

	for _, edge := range edges {
		_, err := repo.Edges().SaveEdge(edge)
		require.NoError(t, err)
	}

	return repo
}

func requireQueryEdgeNames(t *testing.T, edges []*nod.Edge, expected ...string) {
	t.Helper()

	actual := make([]string, 0, len(edges))
	for _, edge := range edges {
		actual = append(actual, edge.Core.Name)
	}

	sort.Strings(actual)
	sort.Strings(expected)
	require.Equal(t, expected, actual)
}
