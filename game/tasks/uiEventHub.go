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

type HandlerInfo struct {
	Id      int
	Handler Handler
}

type Handler func(Event)

type EventHub struct {
	subscribers map[string][]HandlerInfo
}

func NewEventHub() *EventHub {
	return &EventHub{
		subscribers: make(map[string][]HandlerInfo),
	}
}

func (h *EventHub) Subscribe(eventType Event, handler Handler) int {
	eventId++
	t := reflect.TypeOf(eventType).String()
	h.subscribers[t] = append(h.subscribers[t], HandlerInfo{
		Id:      eventId,
		Handler: handler,
	})
	return eventId
}

func (h *EventHub) Publish(event Event) {
	fmt.Printf("Publishg: %T\n", event)
	t := reflect.TypeOf(event).String()
	for _, handlerInfo := range h.subscribers[t] {
		handlerInfo.Handler(event)
	}
}

func (h *EventHub) Unsubscribe(eventType Event, id int) {
	t := reflect.TypeOf(eventType).String()
	handlers := h.subscribers[t]

	for i, info := range handlers {
		if info.Id == id {
			h.subscribers[t] = append(handlers[:i], handlers[i+1:]...)
			return
		}
	}
}
