package entities

type StateMachine struct {
	EveryUpdate                      []EntityUpdater
	EveryUpdateEarlyReturnConditions []func(ent *Entity, gs GameState) bool
	stateController                  StateController // can be nil, maps next state if we want a non-sequential transition between states
	States                           map[int]*StateHandler
	CurrentState                     int
	ResetFunc                        EntityTransitioner
}

type StateController func(entity *Entity) int
type EntityTransitioner func(entity *Entity)
type EntityUpdater func(entity *Entity, gs GameState)

func (s *StateMachine) Reset(ent *Entity) {
	if s.ResetFunc != nil {
		s.ResetFunc(ent)
	}
	s.CurrentState = 1
}

func (s *StateMachine) Transition(ent *Entity) {

	if s.States[s.CurrentState].TransitionOutFunc != nil {
		s.States[s.CurrentState].TransitionOutFunc(ent)
	}

	if s.stateController != nil {
		s.CurrentState = s.stateController(ent)
		s.States[s.CurrentState].TransitionToFunc(ent)
		return
	}

	s.CurrentState = s.States[s.CurrentState].TransitionTo

	if s.States[s.CurrentState] == nil {
		s.Reset(ent)
	}

}

func (ent *Entity) Transition() {
	ent.StateMachine.Transition(ent)
}

type StateHandler struct {
	Updater           func(entity *Entity, state GameState)
	TransitionTo      int
	TransitionOutFunc func(entity *Entity)
	TransitionToFunc  func(entity *Entity)
}

func (s *StateMachine) Update(ent *Entity, gs GameState) {

	//early return(dont update at all) based on every update return condition
	//aim to not modify state with these.
	//Should be something like freeze timer but state is maintained while timer is active, so when we unfreeze ent is in the previous state.
	//NOT timer that freezes, modifies state when it's done then unfreeze, state modification should come from state updater sequence or
	//state controller function only.
	for _, up := range s.EveryUpdateEarlyReturnConditions {
		if up(ent, gs) {
			return
		}
	}

	//some entities have variables states with universal behaviour, just add the func here instead of adding it to each function
	//we could set a flag to not do the common updater if necessary for edge cases
	for _, up := range s.EveryUpdate {
		up(ent, gs)
	}

	s.States[s.CurrentState].Updater(ent, gs)
}
