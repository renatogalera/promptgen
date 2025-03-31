// File: pkg/clipboard/clipboard.go
package clipboard

import (
	"github.com/atotto/clipboard"
)

// Manager provides clipboard operations
type Manager struct{}

// New creates a new clipboard manager
func New() *Manager {
	return &Manager{}
}

// Copy copies text to the system clipboard
func (m *Manager) Copy(text string) error {
	return clipboard.WriteAll(text)
}

// Paste gets text from the system clipboard
func (m *Manager) Paste() (string, error) {
	return clipboard.ReadAll()
}