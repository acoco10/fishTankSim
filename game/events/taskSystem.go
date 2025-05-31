package events

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

func (t *Task) CheckForCompletion() {
	if t.CurrentCount >= t.CountRequired {
		t.Completed = true
	}
}

func (t *Task) PublishIfCompleted(hub *EventHub) {
	if t.Completed == true {
		ev := TaskRequirementsCompleted{
			Task: *t,
		}
		hub.Publish(ev)
	}
}

func (t *Task) Publish(hub *EventHub) {
	ev := TaskCreated{
		Task: *t,
	}
	hub.Publish(ev)
}

func (t *Task) Subscribe(hub *EventHub) {
	hub.Subscribe(t.EventType, func(e Event) {
		if t.Condition == nil || t.Condition(e) && !t.Completed {
			t.CurrentCount++
			t.CheckForCompletion()
			t.PublishIfCompleted(hub)
		}
	})
}
