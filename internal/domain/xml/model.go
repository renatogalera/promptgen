package xml

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"github.com/renatogalera/promptgen/internal/domain/prompt"
)

type CDATA string

func (c CDATA) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(struct {
		Data string `xml:",cdata"`
	}{Data: string(c)}, start)
}

type XMLPrompt struct {
	XMLName     xml.Name `xml:"prompt"`
	Title       string   `xml:"title"`
	Tags        string   `xml:"tags,omitempty"`
	Description string   `xml:"description,omitempty"`
	Content     CDATA    `xml:"content"`
}

type Formatter struct{}

func NewFormatter() *Formatter {
	return &Formatter{}
}

func (f *Formatter) FormatAsXML(p prompt.Prompt) (string, error) {
	x := XMLPrompt{
		Title:       p.Title,
		Tags:        strings.Join(p.Tags, ", "),
		Description: p.Description,
		Content:     CDATA(p.Content),
	}

	data, err := xml.MarshalIndent(x, "", "    ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal prompt to XML: %w", err)
	}

	return xml.Header + string(data), nil
}

func (f *Formatter) FormatWithDoc(p prompt.Prompt) (string, error) {
	baseOutput, err := f.FormatAsXML(p)
	if err != nil {
		return "", err
	}

	if p.Doc == "" {
		return baseOutput, nil
	}

	docContent, err := os.ReadFile(p.Doc)
	if err != nil {
		return baseOutput, fmt.Errorf("error loading doc file '%s': %w", p.Doc, err)

	}

	fullOutput := baseOutput + "\n\n<doc>\n" + string(docContent) + "\n</doc>"
	return fullOutput, nil
}
