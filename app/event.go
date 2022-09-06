package app

const (
	EventAddToUploadList = "eventAddToUploadList"
	EventBackend         = "eventBackend"

	TypeErrorEvent            = "error"
	TypeListObjectsEvent      = "listObjects"
	TypeShowTasksEvent        = "showTasks"
	TypeListenUploadTaskEvent = "listenUploadTask"
)

type Event struct {
	Type string   `json:"type"`
	Args []string `json:"args"`
}
