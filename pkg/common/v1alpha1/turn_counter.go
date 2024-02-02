package v1alpha1

type TurnCounter struct {
	Players      int
	CurrentIndex int
}

func NewTurnCounter(players int, start int) TurnCounter {
	return TurnCounter{
		Players:      players,
		CurrentIndex: start,
	}
}

func (tc TurnCounter) Next() int {
	return (tc.CurrentIndex + 1) % tc.Players
}

func (tc TurnCounter) CurrentPlayer() int {
	return tc.CurrentIndex
}

func (tc TurnCounter) Advance() TurnCounter {
	return TurnCounter{
		Players:      tc.Players,
		CurrentIndex: tc.Next(),
	}
}
