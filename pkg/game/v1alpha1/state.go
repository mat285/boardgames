package v1alpha1

type State struct {
	Version  uint64
	Players  []Player
	Attempts int
	Data     StateData
}

func NewState(players []Player) *State {
	return &State{
		Version: 0,
		Players: players,
	}
}
