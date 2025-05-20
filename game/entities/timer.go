package entities

type TimerState uint8

const (
	Reset TimerState = iota
	Active
	Done
)

type Timer struct {
	Duration int //ok as int because ebiten only uses ticks(no fractional tics)(I think)
	Elapsed  int
	TimerState
	on bool
}

func NewTimer(durationSeconds float64) *Timer {
	//human convenience if fractional seconds to ticks are annoying switch back to ticks as input
	durationTicks := int(durationSeconds * 60)
	t := Timer{}
	t.TimerState = Reset
	t.Duration = durationTicks
	t.Elapsed = 0
	t.on = false
	return &t
}

func (t *Timer) Update() TimerState {

	if t.on {
		switch t.TimerState {
		case Active:
			t.Elapsed += 1
		case Done:
			t.Reset()
		case Reset:
			t.TimerState = Active
		}
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

func (t *Timer) TurnOn() {
	t.on = true
}

func (t *Timer) TurnOff() {
	t.Reset()
	t.on = false
	t.TimerState = Reset
}
