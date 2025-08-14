package tasks

import (
	"fmt"
	"reflect"
)

var eventId int
var EventMap = make(map[int]Handler)

type Event interface {
	Type() string
}

type Handler func(Event)

type EventHub struct {
	subscribers map[string][]Handler
}

func NewEventHub() *EventHub {
	return &EventHub{
		subscribers: make(map[string][]Handler),
	}
}

func (h *EventHub) Subscribe(eventType Event, handler Handler) {
	fmt.Printf("Subscribed to event: %T\n", eventType)
	t := reflect.TypeOf(eventType).String()
	h.subscribers[t] = append(h.subscribers[t], handler)
}

func (h *EventHub) Publish(event Event) {
	fmt.Printf("Publishg: %T\n", event)
	t := reflect.TypeOf(event).String()
	for _, handler := range h.subscribers[t] {
		handler(event)
	}
}
