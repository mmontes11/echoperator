package controller

import echov1alpha1 "github.com/mmontes11/echoperator/pkg/echo/v1alpha1"

type eventType string

const (
	add    eventType = "add"
	delete eventType = "delete"
	update eventType = "update"
)

type event struct {
	key       string
	eventType eventType
	newEcho   *echov1alpha1.Echo
}
