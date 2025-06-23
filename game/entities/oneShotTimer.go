package entities

type OneShotTimer struct {
	*Timer
}

func (t *OneShotTimer) Update() TimerState {

	if t.on {
		switch t.TimerState {
		case Active:
			t.Elapsed += 1
		}
	}

	if t.Duration == t.Elapsed {
		t.TimerState = Done
	}

	return t.TimerState
}
