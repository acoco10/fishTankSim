package events

type TaskRequirementsCompleted struct {
	Task Task
}

type TaskCompleted struct {
	Task Task
}

type TaskCreated struct {
	Task Task
}
