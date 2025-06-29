package lsp

// TODO: support file operations
type Workspace struct {
	WorkspaceFolders WorkspaceFoldersServerCapabilities `json:"workspaceFolders"`
}

type WorkspaceFoldersServerCapabilities struct {
	Supported           bool `json:"supported"`
	ChangeNotifications bool `json:"changeNotifications"`
}

type WorkspaceFolder struct {
	URI  string `json:"uri"`
	Name string `json:"name"`
}

type DidChangeWorkspaceFoldersNotifications struct {
	Notification
	Params DidChangeWorkspaceFoldersParams `json:"params"`
}

type DidChangeWorkspaceFoldersParams struct {
	Event WorkspaceFoldersChangeEvent `json:"event"`
}

type WorkspaceFoldersChangeEvent struct {
	Added   []WorkspaceFolder `json:"added"`
	Removed []WorkspaceFolder `json:"removed"`
}

type DidChangeWatchedFilesNotification struct {
	Notification
	Params DidChangeWatchedFilesParams `json:"params"`
}

type DidChangeWatchedFilesParams struct {
	Changes []FileEvent `json:"changes"`
}

type FileEvent struct {
	URI  string         `json:"uri"`
	Type FileChangeType `json:"type"`
}

type FileChangeType int

const (
	Created FileChangeType = 1
	Changed FileChangeType = 2
	Deleted FileChangeType = 3
)
