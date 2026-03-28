package nod

type Manager struct {
	repository *Repository
}

func NewManager(repository *Repository) *Manager {
	return &Manager{
		repository: repository,
	}
}

func (m *Manager) ExecuteE(fn func(repository *Repository) error) error {
	return m.repository.Transaction(fn)
}

func (m *Manager) Execute(fn func(repository *Repository) error) {
	err := m.ExecuteE(fn)
	if err != nil {
		m.repository.Log.Error("Error executing function in transaction", "error", err)
	}
}
