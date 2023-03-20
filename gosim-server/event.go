package main

type event struct {
	EventType string
	Data      interface{}
}

// Event Types
type chatMessageEvent struct {
	FromUser string
	ToUsers  []string
	Message  string
}

func NewEvent(eventType string, data interface{}) event {
	return event{
		EventType: eventType,
		Data:      data,
	}
}
