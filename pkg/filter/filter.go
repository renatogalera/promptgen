// File: pkg/filter/filter.go
package filter

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
)

// MultiTokenSubstringFilter checks if *all* search tokens (lowercased) are substrings
// of the target (also lowercased). Returns a list of matching items by index.
func MultiTokenSubstringFilter(term string, targets []string) []list.Rank {
	searchTokens := strings.Fields(strings.ToLower(term))
	var ranks []list.Rank

	for i, target := range targets {
		t := strings.ToLower(target)
		matchesAll := true

		for _, token := range searchTokens {
			if !strings.Contains(t, token) {
				matchesAll = false
				break
			}
		}
		if matchesAll {
			// Use default Rank fields as we don't need custom sorting logic here
			ranks = append(ranks, list.Rank{Index: i})
		}
	}
	return ranks
}