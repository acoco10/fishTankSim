package gameEntities

type Task struct {
	Text      string
	hub       *EventHub
	name      string
	completed bool
}

func NewTask(text string, hub *EventHub) *Task {
	t := Task{}
	t.Text = text
	t.hub = hub
	t.completed = false
	t.subs()
	return &t
}

func (t *Task) subs() {
	t.hub.Subscribe(CreatureReachedPoint{}, func(e Event) {
		ev := e.(CreatureReachedPoint)
		if ev.Point.PType == Food {
			if !t.completed {
				evSend := TaskRequirementsCompleted{}
				evSend.task = t.Text
				t.hub.Publish(evSend)
				t.completed = true
			}
		}
	})

	ev := SendData{
		DataFor: "whiteBoard",
		Data:    t.Text,
	}

	t.hub.Publish(ev)
}
