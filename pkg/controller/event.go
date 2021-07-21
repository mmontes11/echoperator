package controller

type eventType string

const (
	addEcho          eventType = "addEcho"
	addScheduledEcho eventType = "addScheduledEcho"
)

type event struct {
	eventType eventType
	resource  interface{}
}
