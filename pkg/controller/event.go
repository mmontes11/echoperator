package controller

type eventType string

const (
	addEcho             eventType = "addEcho"
	addScheduledEcho    eventType = "addScheduledEcho"
	updateScheduledEcho eventType = "updateScheduledEcho"
)

type event struct {
	eventType      eventType
	oldObj, newObj interface{}
}
