package game

import (
	"encoding/json"
	"fmt"

	"github.com/blend/go-sdk/uuid"
	"github.com/mat285/boardgames/games/splendor/pkg/items"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
)

var (
	_ v1alpha1.StateData = State{}
)

type State struct {
	Config       Config
	Players      []Player
	CurrentIndex int

	Board items.Board
}

func NewState(players []Player, config Config) State {
	return State{
		Players:      players,
		CurrentIndex: 0,
		Config:       config,
		Board:        items.NewBoard(),
	}
}

func (s State) apply(move Move) (state State, valid bool, err error) {
	if move.Collect != nil {
		state, valid, err = s.applyCollect(*move.Collect)
	} else if move.Purchase != nil {
		state, valid, err = s.applyPurchase(*move.Purchase)
	} else if move.Reserve != nil {
		state, valid, err = s.applyReserve(*move.Reserve)
	} else if move.Pass != nil {
		state, valid, err = s, true, nil
	} else {
		return s, false, fmt.Errorf("No move")
	}
	if valid {
		state.CurrentIndex = s.Next()
	}
	return
}

func (s State) applyCollect(move CollectMove) (State, bool, error) {
	err := move.Validate(s.Board.Gems)
	if err != nil {
		return s, false, err
	}
	hand := s.Players[s.CurrentIndex].Hand
	gems := hand.Gems.Add(move.Take).ToMap()

	if gems.Total() > 10 {
		gems.Sub(move.Return.ToMap())
		if !gems.NonNegative() || gems.Total() > 10 {
			return s, false, err
		}
	}
	hand.Gems = gems.ToCount()
	s.Players[s.CurrentIndex].Hand = hand
	return s, true, err
}

func (s State) applyPurchase(move CardMove) (State, bool, error) {
	hand := s.Players[s.CurrentIndex].Hand

	if !hand.CanPurchase(move.Card) {
		return s, false, nil
	}
	reserved := items.ContainsCard(hand.Reserved, move.Card)
	onBoard := s.Board.IsCardOnBoard(move.Card)

	if !reserved && !onBoard {
		return s, false, nil
	}

	hand = hand.Purchase(move.Card)
	bonuses := hand.BonusesEarned(s.Board.Bonuses)
	hand.Bonus = append(hand.Bonus, bonuses...)
	s.Board = s.Board.RemoveBonuses(bonuses)
	if onBoard {
		s.Board = s.Board.RemoveCard(move.Card)
	}
	s.Players[s.CurrentIndex].Hand = hand
	return s, true, nil
}

func (s State) applyReserve(move CardMove) (State, bool, error) {
	hand := s.Players[s.CurrentIndex].Hand
	if !hand.CanReserve() {
		return s, false, nil
	}
	if !s.Board.IsCardOnBoard(move.Card) {
		return s, false, nil
	}
	s.Board = s.Board.RemoveCard(move.Card)
	hand = hand.Reserve(move.Card)
	s.Players[s.CurrentIndex].Hand = hand
	return s, true, nil
}

func (s State) CurrentPlayer() (uuid.UUID, error) {
	if len(s.Players) <= s.CurrentIndex {
		return nil, fmt.Errorf("Invalid Player Index %d", s.CurrentIndex)
	}
	return s.Players[s.CurrentIndex].ID, nil
}

func (s State) IsDone() bool {
	return len(s.Winners()) > 0
}

func (s State) Next() int {
	return (s.CurrentIndex + 1) % len(s.Players)
}

func (s State) AdvancePlayer() {
	s.CurrentIndex = s.Next()
}

func (s State) Winners() []uuid.UUID {
	maxP := 0
	var max []uuid.UUID
	for _, p := range s.Players {
		pts := p.Hand.Points()
		if pts >= s.Config.VictoryPoints {
			if maxP < pts {
				max = []uuid.UUID{p.ID}
				maxP = pts
			} else if maxP == pts {
				max = append(max, p.ID)
			}
		}
	}
	return max
}

func (s State) Serialize() ([]byte, error) {
	hand, err := json.MarshalIndent(s.Players[s.CurrentIndex].Hand, "", "  ")
	if err != nil {
		return nil, err
	}
	cards, err := json.MarshalIndent(s.Board.AvailableCards(), "", "  ")
	if err != nil {
		return nil, err
	}
	hand = append(hand, '\n', '\n')
	return append(hand, cards...), nil
}

func (s State) Deserialize(data []byte) error {
	return json.Unmarshal(data, &s)
}
