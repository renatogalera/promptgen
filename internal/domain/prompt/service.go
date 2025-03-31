package prompt

import (
	"sort"
	"strings"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) SortPromptsByTitle(collection *PromptCollection) {
	sort.SliceStable(collection.Prompts, func(i, j int) bool {
		return strings.ToLower(collection.Prompts[i].Title) < strings.ToLower(collection.Prompts[j].Title)
	})
}

func (s *Service) ParseTags(tagString string) []string {
	if strings.TrimSpace(tagString) == "" {
		return nil
	}

	tags := strings.Split(tagString, ",")
	cleanedTags := make([]string, 0, len(tags))

	for _, t := range tags {
		trimmed := strings.TrimSpace(t)
		if trimmed != "" {
			cleanedTags = append(cleanedTags, trimmed)
		}
	}

	if len(cleanedTags) == 0 {
		return nil
	}

	return cleanedTags
}

func (s *Service) ReplaceVariables(content string, vars map[string]string) string {
	output := content
	for k, v := range vars {
		placeholder := "{{{" + k + "}}}"
		output = strings.ReplaceAll(output, placeholder, v)
	}
	return output
}

func (s *Service) CreatePrompt(title, tagString, description, content string) Prompt {
	return Prompt{
		Title:       title,
		Tags:        s.ParseTags(tagString),
		Description: description,
		Content:     content,
	}
}

func (s *Service) AddPrompt(collection *PromptCollection, prompt Prompt) {
	collection.Prompts = append(collection.Prompts, prompt)
}
