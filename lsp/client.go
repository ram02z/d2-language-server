package lsp

import "github.com/google/uuid"

type RegistrationParams struct {
	Registrations []Registration `json:"registrations"`
}

type Registration struct {
	ID              string `json:"id"`
	Method          string `json:"method"`
	RegisterOptions any    `json:"registerOptions"`
}

type DidChangeWatchedFilesRegistrationOptions struct {
	Watchers []FileSystemWatcher `json:"watchers"`
}

type FileSystemWatcher struct {
	GlobPattern string `json:"globPattern"`
}

func NewRegistration(method Method, registerOptions any) Registration {
	return Registration{
		ID:              uuid.New().String(),
		Method:          string(method),
		RegisterOptions: registerOptions,
	}
}
