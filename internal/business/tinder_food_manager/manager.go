package tinderfood

import (
	"fmt"
)

type Manager struct {
	repo repo
}

func NewManager(repo repo) *Manager {
	return &Manager{repo: repo}
}

func (m *Manager) CreateSession() (string, error) {
	session, err := m.repo.CreateSession(nil)
	if err != nil {
		return "", fmt.Errorf("failed to create session err:%w", err)
	}

	return session.ID.String(), nil
}
