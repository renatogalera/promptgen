package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/renatogalera/promptgen/internal/app"
	"github.com/renatogalera/promptgen/internal/config"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "promptgen[flags]",
	Short: "Interactive AI prompt management system in your terminal",
	Long: `promptgenhelps you organize, search, and use your AI prompts efficiently.
Load prompts from a YAML file, browse, search (press '/'), view details,
and copy prompts formatted as XML to the clipboard. Create new prompts using 'n'.`,
	Run: func(cmd *cobra.Command, args []string) {

		promptFile, _ := cmd.Flags().GetString("file")

		promptFile = resolvePromptFilePath(promptFile)

		application := app.NewApplication(promptFile)
		p := tea.NewProgram(application, tea.WithAltScreen(), tea.WithMouseCellMotion())

		if _, err := p.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
			os.Exit(1)
		}
	},
}

func resolvePromptFilePath(promptFile string) string {

	if promptFile != "" {
		fmt.Fprintf(os.Stderr, "Info: Using specified prompts file: %s\n", promptFile)
		ensureDirectoryExists(filepath.Dir(promptFile))
		return promptFile
	}

	if _, err := os.Stat(config.DefaultPromptFilename); err == nil {
		promptFile = config.DefaultPromptFilename
		fmt.Fprintf(os.Stderr, "Info: Using local prompts file: %s\n", promptFile)
		return promptFile
	}

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error: Could not resolve home directory: %v\n", err)
	}

	configDir := filepath.Join(home, ".config", config.ConfigDirName)
	configPath := filepath.Join(configDir, config.DefaultPromptFilename)

	if _, err := os.Stat(configPath); err == nil {
		promptFile = configPath
		fmt.Fprintf(os.Stderr, "Info: Using config prompts file: %s\n", promptFile)
	} else {
		promptFile = configPath
		fmt.Fprintf(os.Stderr, "Info: No prompts file found. Will create/use: %s\n", promptFile)
		ensureDirectoryExists(configDir)
	}

	return promptFile
}

func ensureDirectoryExists(dir string) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("Error: Could not create directory '%s': %v\n", dir, err)
	}
}

func init() {
	defaultPathDesc := fmt.Sprintf("./%s or ~/.config/%s/%s",
		config.DefaultPromptFilename,
		config.ConfigDirName,
		config.DefaultPromptFilename)

	rootCmd.Flags().StringP("file", "f", "", "Path to the prompts YAML file (default: "+defaultPathDesc+")")
}
