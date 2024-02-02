package game

import (
	"fmt"

	"github.com/blend/go-sdk/uuid"
	"github.com/mat285/boardgames/games/splendor/meta"
	"github.com/mat285/boardgames/games/splendor/pkg/items"
	common "github.com/mat285/boardgames/pkg/common/v1alpha1"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
)

var (
	_ v1alpha1.StateData = State{}
)

type State struct {
	meta.Object
	Config  Config
	Players []Player

	Turn common.TurnCounter

	Board items.Board
}

func NewState(players []Player, config Config) State {
	return State{
		Players: players,
		Config:  config,
		Turn:    common.NewTurnCounter(len(players), 0),
		Board:   items.NewBoard(),
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
		state.Turn = state.Turn.Advance()
	}
	return
}

func (s State) applyCollect(move CollectMove) (State, bool, error) {
	err := move.Validate(s.Board.Gems)
	if err != nil {
		return s, false, err
	}
	player, err := s.GetCurrentPlayer()
	if err != nil {
		return s, false, err
	}
	hand := player.Hand
	gems := hand.Gems.Add(move.Take).ToMap()

	if gems.Total() > 10 {
		gems.Sub(move.Return.ToMap())
		if !gems.NonNegative() || gems.Total() > 10 {
			return s, false, err
		}
	}
	hand.Gems = gems.ToCount()
	s = s.setCurrentPlayerHand(hand)
	return s, true, err
}

func (s State) applyPurchase(move CardMove) (State, bool, error) {
	player, err := s.GetCurrentPlayer()
	if err != nil {
		return s, false, err
	}
	hand := player.Hand
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

	s = s.setCurrentPlayerHand(hand)
	return s, true, nil
}

func (s State) applyReserve(move CardMove) (State, bool, error) {
	player, err := s.GetCurrentPlayer()
	if err != nil {
		return s, false, err
	}
	hand := player.Hand
	if !hand.CanReserve() {
		return s, false, nil
	}
	if !s.Board.IsCardOnBoard(move.Card) {
		return s, false, nil
	}
	s.Board = s.Board.RemoveCard(move.Card)
	hand = hand.Reserve(move.Card)
	player.Hand = hand
	return s, true, nil
}

func (s State) CurrentPlayer() (uuid.UUID, error) {
	player, err := s.GetCurrentPlayer()
	if err != nil {
		return nil, err
	}
	return player.ID, nil
}

func (s State) GetCurrentPlayer() (Player, error) {
	idx := s.Turn.CurrentPlayer()
	if len(s.Players) <= idx {
		return Player{}, fmt.Errorf("Invalid Player Index %d", idx)
	}
	return s.Players[idx], nil
}

func (s State) setCurrentPlayer(p Player) State {
	idx := s.Turn.CurrentPlayer()
	if len(s.Players) > idx {
		s.Players[idx] = p
	}
	return s
}

func (s State) setCurrentPlayerHand(h items.Hand) State {
	idx := s.Turn.CurrentPlayer()
	if len(s.Players) > idx {
		s.Players[idx].Hand = h
	}
	return s
}

func (s State) IsDone() bool {
	return len(s.Winners()) > 0
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
