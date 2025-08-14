package tasks

type TaskRequirementsCompleted struct {
	Task Task
}

func (t TaskRequirementsCompleted) Type() string {
	return "TaskRequirementsCompleted"
}

type TaskCompleted struct {
	Task Task
	Slot int
}

func (t TaskCompleted) Type() string {
	return "TaskCompleted"
}

type TaskCreated struct {
	Task *Task
}

func (t TaskCreated) Type() string {
	return "TaskCreated"
}

type AllTasksCompleted struct {
}

func (a AllTasksCompleted) Type() string {
	return "AllTasksCompleted"
}
