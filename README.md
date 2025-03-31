# PromptGen

![Go Version](https://img.shields.io/badge/Go-1.16%2B-blue)

**PromptGen** is an interactive terminal application for efficiently managing, organizing, and utilizing AI prompts. The system allows you to load prompts from a YAML file, browse, search, and copy prompts formatted as XML to the clipboard.

## 📋 Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)

## 📝 Introduction

PromptGen is a tool designed for professionals and enthusiasts working with AI models who need to manage a collection of prompts. The application provides a user-friendly terminal interface (TUI) built with [Bubble Tea](https://github.com/charmbracelet/bubbletea), offering:

- Centralized storage for all your AI prompts
- Efficient organization with titles, tags, and descriptions
- Fast searching and filtering capabilities
- Simplified access and usage with clipboard integration
- Support for variable placeholders in prompts

## ✨ Features

- **Intuitive Navigation**: Browse your prompt library using keyboard shortcuts
- **Powerful Search**: Quickly find prompts with multi-token search functionality
- **Rich Prompt Details**: View comprehensive information about each prompt
- **XML Formatting**: Copy prompts as properly formatted XML to the clipboard (inspired by Anthropic's recommended approach for structuring prompts)
- **Variable Support**: Customize prompts with variable placeholders
- **Prompt Creation**: Add new prompts directly from the interface
- **Local Storage**: Manage prompts using a simple YAML file format

## 🛠️ Requirements

- Go 1.16 or higher
- Operating system compatible with the clipboard library (Linux, macOS, Windows)
- Dependencies (automatically installed via Go modules):
  - github.com/charmbracelet/bubbles
  - github.com/charmbracelet/bubbletea
  - github.com/charmbracelet/lipgloss
  - github.com/atotto/clipboard
  - github.com/spf13/cobra
  - gopkg.in/yaml.v3

## 📦 Installation

### Via Go Install

```bash
go install github.com/renatogalera/promptgen/cmd/promptgent@latest
```

### Building from Source

```bash
# Clone the repository
git clone https://github.com/renatogalera/promptgen.git
cd promptgen

# Build the project
go build -o promptgen./cmd/promptgent

# Optional: move the binary to a directory in your PATH
sudo mv promptgen/usr/local/bin/
```

## 🚀 Usage

Run the program without arguments to start with default settings:

```bash
promptgent
```

PromptGen will look for a prompts file in:
1. Current directory: `./prompts.yaml`
2. Configuration directory: `~/.config/promptgen/prompts.yaml`

If no file is found, a new one will be automatically created.

### Keyboard Shortcuts

| Key | Function |
|-----|----------|
| `/` | Search prompts |
| `↑/↓` or `j/k` | Navigate the list |
| `Enter` | Select prompt |
| `c` | Copy prompt as XML to clipboard |
| `d` | Copy associated documentation (if available) |
| `n` | Create new prompt |
| `Esc` | Go back |
| `?` | Show help |
| `Ctrl+C` | Exit application |
| `Tab/Shift+Tab` | Navigate form fields when creating or editing prompts |

## ⚙️ Configuration

### Custom Prompts File

To use a custom prompts file:

```bash
promptgen-f /path/to/my_prompts.yaml
```

### YAML File Structure

```yaml
prompts:
  - title: "My Prompt"
    tags: ["tag1", "tag2"]
    description: "Description of the prompt"
    content: "Content of the prompt with {{{variable1}}} placeholders"
    variables: ["variable1", "variable2"]
    doc: "/optional/path/to/documentation.txt"  # Documentation file to be included with the prompt
  
  - title: "Another Prompt"
    # ...
```

The `doc` field allows you to specify a path to a documentation file that will be included with the prompt when copied. This documentation can be accessed separately using the `d` key in the prompt view.

## 📋 Examples

### Example Workflow

1. Start PromptGen: `promptgent`
2. Navigate through prompts using arrow keys or j/k
3. Press `/` to search for a specific prompt (supports multi-token search)
4. Select a prompt with `Enter` to view details
5. Press `c` to copy the formatted XML prompt to the clipboard
6. For prompts with variables, press `Enter` again to input values
7. Create a new prompt by pressing `n` in the list screen

### Variable Placeholders

You can define variables in your prompts using the `{{{variableName}}}` syntax. When using a prompt with variables, PromptGen will prompt you to enter values for each variable before generating the final output.

### XML Output Example

The XML format used for prompts was inspired by [Anthropic's recommendation](https://docs.anthropic.com/en/docs/build-with-claude/prompt-engineering/use-xml-tags) to use XML tags when working with AI assistants like Claude. This structured format helps AI models better understand and process different parts of your prompts.

```xml
<?xml version="1.0" encoding="utf-8"?>
<prompt>
  <title>My Prompt</title>
  <tags>tag1, tag2</tags>
  <description>Description of the prompt</description>
  <content><![CDATA[Content of the prompt with inserted variable values]]></content>
</prompt>

<doc>
Additional documentation content if a doc file is specified in the prompt configuration
</doc>
```

The `doc` element will only be included when the prompt has an associated documentation file specified in the `doc` field of the YAML configuration.


## 👥 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the project
2. Create a feature branch (`git checkout -b feature/new-feature`)
3. Commit your changes (`git commit -m 'Add new feature'`)
4. Push to the branch (`git push origin feature/new-feature`)
5. Open a Pull Request


## 📜 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

Repository: [github.com/renatogalera/promptgen](https://github.com/renatogalera/promptgen)