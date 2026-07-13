package contract

import (
	"testing"

	"github.com/m87/nod"
)

type RepositoryFactory func(t *testing.T) *nod.Repository

func RunRepositoryContractTests(t *testing.T, factory RepositoryFactory) {
	t.Helper()
}
