package contract

import (
	"testing"

	"github.com/m87/nod"
	"github.com/stretchr/testify/require"
)

func testEdgeQueryCoreFields(t *testing.T, factory RepositoryFactory) {
	repo := createEdgeQueryTestRepository(t, factory)

	tests := []struct {
		name       string
		expression nod.Expression
		expected   []string
	}{
		{
			name:       "id",
			expression: nod.EdgeFields.Id.Equals(queryEdgeAlphaID),
			expected:   []string{"alpha"},
		},
		{
			name:       "name",
			expression: nod.EdgeFields.Name.Equals("beta"),
			expected:   []string{"beta"},
		},
		{
			name:       "namespace id",
			expression: nod.EdgeFields.NamespaceId.Equals(queryNamespaceA),
			expected:   []string{"alpha", "beta"},
		},
		{
			name:       "source id",
			expression: nod.EdgeFields.SourceId.Equals(queryEdgeSourceAID),
			expected:   []string{"alpha", "beta"},
		},
		{
			name:       "target id",
			expression: nod.EdgeFields.TargetId.Equals(queryEdgeTargetAID),
			expected:   []string{"alpha", "gamma"},
		},
		{
			name:       "kind",
			expression: nod.EdgeFields.Kind.Equals("dependency"),
			expected:   []string{"alpha", "beta"},
		},
		{
			name:       "status",
			expression: nod.EdgeFields.Status.Equals("active"),
			expected:   []string{"alpha", "gamma"},
		},
		{
			name:       "in",
			expression: nod.EdgeFields.Kind.In([]string{"dependency", "ownership"}),
			expected:   []string{"alpha", "beta", "delta"},
		},
		{
			name:       "not in",
			expression: nod.EdgeFields.Status.NotIn([]string{"inactive", "archived"}),
			expected:   []string{"alpha", "gamma"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			edges, err := nod.NewEdgeQuery(repo).Where(tt.expression).FindAll()

			require.NoError(t, err)
			requireQueryEdgeNames(t, edges, tt.expected...)
		})
	}
}
