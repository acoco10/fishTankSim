package gameEntities

type TimerState uint8

const (
	Reset TimerState = iota
	Active
	Done
)

type Timer struct {
	Duration int //ok as int because ebiten only uses ticks(no fractional tics)
	Elapsed  int
	TimerState
}

func NewTimer(durationSeconds float64) *Timer {
	//human convenience if fractional seconds to ticks are annoying switch back to ticks as input
	durationTicks := int(durationSeconds * 60)
	t := Timer{}
	t.TimerState = Reset
	t.Duration = durationTicks
	t.Elapsed = 0
	return &t
}

func (t *Timer) Update() TimerState {

	switch t.TimerState {
	case Active:
		t.Elapsed += 1
	case Done:
		t.Reset()
	case Reset:
		t.TimerState = Active
	}

	if t.Duration == t.Elapsed {
		t.TimerState = Done
	}

	return t.TimerState
}

func (t *Timer) Reset() {
	t.Elapsed = 0
	t.TimerState = Reset
}
