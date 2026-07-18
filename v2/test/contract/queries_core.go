package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

func testQueryCoreFields(t *testing.T, factory RepositoryFactory) {
	repo := createQueryTestRepository(t, factory)

	tests := []struct {
		name       string
		expression nod.Expression
		expected   []string
	}{
		{
			name:       "id",
			expression: nod.NodeFields.Id.Equals(queryNodeAlphaID),
			expected:   []string{"alpha"},
		},
		{
			name:       "name",
			expression: nod.NodeFields.Name.Equals("beta"),
			expected:   []string{"beta"},
		},
		{
			name:       "namespace id",
			expression: nod.NodeFields.NamespaceId.Equals(queryNamespaceA),
			expected:   []string{"alpha", "beta"},
		},
		{
			name:       "parent id",
			expression: nod.NodeFields.ParentId.Equals(queryNodeAlphaID),
			expected:   []string{"beta"},
		},
		{
			name:       "kind",
			expression: nod.NodeFields.Kind.Equals("article"),
			expected:   []string{"alpha", "beta"},
		},
		{
			name:       "status",
			expression: nod.NodeFields.Status.Equals("published"),
			expected:   []string{"alpha", "gamma"},
		},
		{
			name:       "in",
			expression: nod.NodeFields.Kind.In([]string{"article", "task"}),
			expected:   []string{"alpha", "beta", "delta"},
		},
		{
			name:       "not in",
			expression: nod.NodeFields.Status.NotIn([]string{"draft", "archived"}),
			expected:   []string{"alpha", "gamma"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes, err := nod.NewNodeQuery(repo).Where(tt.expression).FindAll()

			require.NoError(t, err)
			requireQueryNodeNames(t, nodes, tt.expected...)
		})
	}
}
