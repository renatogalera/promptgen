package yaml

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/renatogalera/promptgen/internal/domain/prompt"
)

type Repository struct {
	filePath string
	service  *prompt.Service
}

func NewRepository(filePath string, service *prompt.Service) *Repository {
	if service == nil {
	}
	return &Repository{
		filePath: filePath,
		service:  service,
	}
}

func (r *Repository) LoadPrompts() (prompt.PromptCollection, error) {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return prompt.PromptCollection{Prompts: []prompt.Prompt{}}, nil
		}
		return prompt.PromptCollection{}, fmt.Errorf("failed to read prompt file '%s': %w", r.filePath, err)
	}
	var pc prompt.PromptCollection
	if err := yaml.Unmarshal(data, &pc); err != nil {
		return prompt.PromptCollection{}, fmt.Errorf("failed to parse YAML from '%s': %w", r.filePath, err)
	}
	r.service.SortPromptsByTitle(&pc)
	return pc, nil
}

func (r *Repository) SavePrompts(collection prompt.PromptCollection) error {
	r.service.SortPromptsByTitle(&collection)
	data, err := yaml.Marshal(collection)
	if err != nil {
		return fmt.Errorf("failed to marshal prompts to YAML: %w", err)
	}
	dir := filepath.Dir(r.filePath)
	if err := os.MkdirAll(dir, 0750); err != nil {
		return fmt.Errorf("failed to create directory '%s': %w", dir, err)
	}
	if err := os.WriteFile(r.filePath, data, 0640); err != nil {
		return fmt.Errorf("failed to save prompt file '%s': %w", r.filePath, err)
	}
	return nil
}

func (r *Repository) SavePrompt(newPrompt prompt.Prompt) error {
	collection, err := r.LoadPrompts()
	if err != nil {

		return fmt.Errorf("could not load existing prompts before saving: %w", err)
	}
	r.service.AddPrompt(&collection, newPrompt)
	if err := r.SavePrompts(collection); err != nil {
		return fmt.Errorf("could not save updated prompts collection: %w", err)
	}
	return nil
}
