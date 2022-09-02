package app

const (
	EventAddToUploadList = "eventAddToUploadList"
	EventBackend         = "eventBackend"

	TypeErrorEvent       = "error"
	TypeListObjectsEvent = "listObjects"
	TypeShowTasksEvent   = "showTasks"
)

type Event struct {
	Type string   `json:"type"`
	Args []string `json:"args"`
}
