package tasks

import (
	"github.com/acoco10/fishTankWebGame/game/events"
)

type EventCondition func(Event) bool

type TaskType uint8

const (
	FishFed TaskType = iota
)

var taskId int
var completedTaskQueue []Event

type Task struct {
	LinkedTask *Task
	Text       string
	Index      int
	Type       TaskType

	Completed  bool
	Condition1 EventCondition
	EventType  Event
	EventType2 Event
	activated  bool

	SubTasks []*SubTask
}

type SubTask struct {
	EventType Event
	Condition EventCondition
	Completed bool
}

func NewTask(EventType Event, text string, condition EventCondition, hub *EventHub) *Task {
	println("creating task id:", taskId)
	task := &Task{
		Text:       text,
		EventType:  EventType,
		Condition1: condition,
		Index:      taskId,
	}
	task.QueueCondition(hub)
	taskId++
	return task
}

func (t *Task) PublishIfCompleted(hub *EventHub) {
	if t.activated {
		ev := TaskRequirementsCompleted{
			Task: *t,
		}
		hub.Publish(ev)
	}
}

func (t *Task) QueueCondition(hub *EventHub) {
	hub.Subscribe(t.EventType, func(e Event) {
		completedTaskQueue = append(completedTaskQueue, e)
	})
}

func (t *Task) Publish(hub *EventHub) {
	ev := TaskCreated{
		Task: t,
	}

	hub.Publish(ev)

}
func (t *Task) Activate(eventHub *EventHub) {
	println("publishing task:", t.Text)
	t.activated = true
	t.Subscribe(eventHub)
	t.Publish(eventHub)

	if CheckCompletedEventQueue(t) {
		ev := TaskRequirementsCompleted{
			Task: *t,
		}
		eventHub.Publish(ev)
	}

	for _, sub := range t.SubTasks {
		sub.Subscribe(eventHub)
	}

	println("publishing task completed after creation after checking queue")

}

func CheckCompletedEventQueue(t *Task) bool {
	for _, event := range completedTaskQueue {
		if event.Type() == t.EventType.Type() {
			if t.Condition1(event) {
				return true
			}
		}
	}
	return false
}

func (t *Task) Activated() bool {
	return t.activated
}

func (t *Task) Subscribe(hub *EventHub) {
	hub.Subscribe(t.EventType, func(e Event) {

		if t.Condition1 == nil || t.Condition1(e) {
			for _, subT := range t.SubTasks {
				if !subT.Completed {
					return
				}
			}
			t.PublishIfCompleted(hub)
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
		ev := e.(events.NewDay)
		if ev.Day > 1 {
			completedTaskQueue = []Event{}
		}
	})
}

func (st *SubTask) Subscribe(hub *EventHub) {

	for _, completedTask := range completedTaskQueue {
		if completedTask.Type() == st.EventType.Type() {
			st.Completed = true
		}
	}

	hub.Subscribe(st.EventType, func(e Event) {
		if st.Condition(e) {
			st.Completed = true
		}
	})
}
