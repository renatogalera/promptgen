package prompt

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
)

type Prompt struct {
	Title       string   `yaml:"title"`
	Tags        []string `yaml:"tags,omitempty"`
	Description string   `yaml:"description,omitempty"`
	Content     string   `yaml:"content"`
	Variables   []string `yaml:"variables,omitempty"`
	Doc         string   `yaml:"doc,omitempty"`
}

type PromptCollection struct {
	Prompts []Prompt `yaml:"prompts"`
}

type Item struct {
	Prompt
}

func (i Item) Title() string {
	return i.Prompt.Title
}

func (i Item) Description() string {
	return strings.Join(i.Tags, ", ")
}

func (i Item) FilterValue() string {
	var builder strings.Builder
	builder.WriteString(i.Prompt.Title)

	if len(i.Tags) > 0 {
		builder.WriteString(" ")
		builder.WriteString(strings.Join(i.Tags, " "))
	}
	return builder.String()
}

func (pc PromptCollection) ToItems() []list.Item {
	items := make([]list.Item, len(pc.Prompts))
	for i, p := range pc.Prompts {
		items[i] = Item{Prompt: p}
	}
	return items
}

func (pc PromptCollection) GetPromptByTitle(title string) (Prompt, bool) {
	for _, p := range pc.Prompts {
		if p.Title == title {
			return p, true
		}
	}
	return Prompt{}, false
}
