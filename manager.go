package nod

// Manager wraps a Repository and provides convenience methods for transactional operations.
type Manager struct {
	repository *Repository
}

// NewManager creates a new Manager wrapping the given Repository.
func NewManager(repository *Repository) *Manager {
	return &Manager{
		repository: repository,
	}
}

// ExecuteE runs fn within a transaction and returns any error.
func (m *Manager) ExecuteE(fn func(repository *Repository) error) error {
	return m.repository.Transaction(fn)
}

// Execute runs fn within a transaction, logging any error.
func (m *Manager) Execute(fn func(repository *Repository) error) {
	err := m.ExecuteE(fn)
	if err != nil {
		m.repository.log.Error("Error executing function in transaction", "error", err)
	}
}
