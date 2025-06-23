package tasks

import "github.com/acoco10/fishTankWebGame/game/events"

type EventCondition func(Event) bool

var taskId int

type Task struct {
	Text          string
	Index         int
	Name          string
	Completed     bool
	CurrentCount  int
	CountRequired int
	Condition     EventCondition
	EventType     Event
	activated     bool
}

func NewTask(EventType Event, text string, condition EventCondition) *Task {
	println("creating task id:", taskId)
	task := &Task{
		Text:          text,
		EventType:     EventType,
		CountRequired: 1,
		Condition:     condition,
		Index:         taskId,
	}
	taskId++
	return task
}

func (t *Task) PublishIfCompleted(hub *EventHub) {
	if t.CurrentCount >= t.CountRequired {
		ev := TaskRequirementsCompleted{
			Task: *t,
		}
		hub.Publish(ev)
	}
}

func (t *Task) Publish(hub *EventHub) {
	ev := TaskCreated{
		Task: t,
	}
	hub.Publish(ev)
}
func (t *Task) Activate() {
	t.activated = true
}

func (t *Task) Activated() bool {
	return t.activated
}

func (t *Task) Subscribe(hub *EventHub) {
	hub.Subscribe(t.EventType, func(e Event) {
		if t.activated {
			if t.Condition == nil || t.Condition(e) && t.CurrentCount < t.CountRequired {
				t.CurrentCount++
				t.PublishIfCompleted(hub)
			}
		}
	})

	hub.Subscribe(TaskCompleted{}, func(e Event) {
		ev := e.(TaskCompleted)
		if ev.Task.Text == t.Text {
			println("received task completed event @ task, updating status for", ev.Task.Text)
			t.Completed = true
		}

	})
	hub.Subscribe(events.NewDay{}, func(e Event) {
		taskId = 0

	})
}
