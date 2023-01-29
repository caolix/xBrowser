package models

const (
	EventAddToUploadList = "eventAddToUploadList"

	AppEventBus = "eventBus"
	// The following events all run over the 'AppEventBus'
	EventListObjects      = "listObjects"
	EventShowTasks        = "showTasks"
	EventListenUploadTask = "listenUploadTask"
)

type Event struct {
	Name string   `json:"name"`
	Args []string `json:"args"`
}
