package v1alpha1

type Event struct {
	Type EventType
	Body interface{}
}

type EventType int

const (
	EventTypeUnknown EventType = 0
	EventTypeStop    EventType = 1
	EventTypeSave    EventType = 2
)
