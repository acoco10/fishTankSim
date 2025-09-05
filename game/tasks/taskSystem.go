package tasks

import "github.com/acoco10/fishTankWebGame/game/events"

type EventCondition func(Event) bool

type TaskType uint8

const (
	FishFed TaskType = iota
)

type TaskManager struct {
	EventHub           *EventHub
	completedTaskQueue []Event
	Tasks              []*Task
	currentTask        int
}

func (tm *TaskManager) Subscribe() {
	tm.EventHub.Subscribe(events.DayOver{}, func(e Event) {
		for _, task := range tm.Tasks {
			tm.EventHub.Unsubscribe(task.EventType, task.Index)
		}
		tm.Tasks = []*Task{}
		tm.completedTaskQueue = []Event{}
	})

	tm.EventHub.Subscribe(TaskCompleted{}, func(e Event) {
		tm.currentTask++
	})

	tm.EventHub.Subscribe(events.NewDay{}, func(e Event) {
		tm.currentTask = 0
	})

}

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

func (tm *TaskManager) NewTask(EventType Event, text string, condition EventCondition) *Task {

	task := &Task{
		Text:       text,
		EventType:  EventType,
		Condition1: condition,
	}

	tm.Tasks = append(tm.Tasks, task)
	tm.QueueCondition(tm.EventHub, *task)
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

func (t *TaskManager) QueueCondition(hub *EventHub, task Task) {
	hub.Subscribe(task.EventType, func(e Event) {
		t.completedTaskQueue = append(t.completedTaskQueue, e)
	})
}

func (t *Task) Publish(hub *EventHub) {
	ev := TaskCreated{
		Task: t,
	}
	hub.Publish(ev)

}

func (t *TaskManager) Activate() {
	task := t.Tasks[t.currentTask]
	println("publishing task:", task.Text)
	task.activated = true
	id := task.Subscribe(t.EventHub)
	task.Index = id
	task.Publish(t.EventHub)

	if t.CheckCompletedEventQueue(*task) {
		ev := TaskRequirementsCompleted{
			Task: *task,
		}
		t.EventHub.Publish(ev)
	}

	println("publishing task completed after creation after checking queue")

}

func (tm *TaskManager) CheckCompletedEventQueue(task Task) bool {
	for _, event := range tm.completedTaskQueue {
		if event.Type() == task.EventType.Type() {
			if task.Condition1(event) {
				return true
			}
		}
	}
	return false
}

func (t *Task) Activated() bool {
	return t.activated
}

func (t *Task) Subscribe(hub *EventHub) int {
	id := hub.Subscribe(t.EventType, func(e Event) {
		if t.Condition1 == nil || t.Condition1(e) {
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

	return id
}
