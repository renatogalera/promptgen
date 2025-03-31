package config

import (
	"time"
)

const (
	DefaultPromptFilename = "prompts.yaml"
	ConfigDirName         = "promptgen"
	StatusMessageTimeout  = 3 * time.Second
)

type AppState int

const (
	StateLoading AppState = iota
	StatePromptList
	StatePromptView
	StatePromptCreation
	StateVariableInput
)
