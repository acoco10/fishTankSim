package tasks

type TaskRequirementsCompleted struct {
	Task Task
}

type TaskCompleted struct {
	Task Task
	Slot int
}

type TaskCreated struct {
	Task *Task
}

type AllTasksCompleted struct {
}
