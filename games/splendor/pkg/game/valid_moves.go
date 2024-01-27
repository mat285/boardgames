package game

import (
	"github.com/mat285/boardgames/games/splendor/pkg/items"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
)

func (s State) ValidMoves() ([]v1alpha1.Move, error) {
	moves := s.ValidCollectMoves()
	moves = append(moves, s.ValidPurchaseMoves()...)
	moves = append(moves, s.ValidReserveMoves()...)
	return MoveSliceToMoveSlice(moves), nil
}

func (s State) ValidCollectMoves() []*Move {
	moves := []*Move{}
	hand := s.Players[s.CurrentIndex].Hand
	if hand.Gems.Total() > 8 {
		return moves
	}

	gems := s.Board.Gems.ToMap()
	delete(gems, items.GemWild)
	sets := allSetsOf3(gems.Keys())

	for _, set := range sets {
		moves = append(moves, &Move{Collect: &CollectMove{Take: set}})
	}

	for k, v := range gems {
		if v >= 4 {
			count := items.GemMap{k: 2}
			moves = append(moves, &Move{Collect: &CollectMove{Take: count.ToCount()}})
		}
	}
	return moves
}

func (s State) ValidPurchaseMoves() []*Move {
	moves := []*Move{}
	hand := s.Players[s.CurrentIndex].Hand
	for _, card := range s.Board.AvailableCards() {
		if hand.CanPurchase(card) {
			moves = append(moves, &Move{Purchase: &CardMove{Card: card}})
		}
	}
	for _, card := range hand.Reserved {
		if hand.CanPurchase(card) {
			moves = append(moves, &Move{Purchase: &CardMove{Card: card}})
		}
	}
	return moves
}

func (s State) ValidReserveMoves() []*Move {
	moves := []*Move{}
	hand := s.Players[s.CurrentIndex].Hand
	if !hand.CanReserve() || s.Board.Gems.Wild <= 0 {
		return moves
	}
	for _, card := range s.Board.AvailableCards() {
		moves = append(moves, &Move{Reserve: &CardMove{Card: card}})
	}
	return moves
}

func allSetsOf3(gems []items.Gem) []items.GemCount {
	sets := make([]items.GemCount, 0)
	for i := 0; i < len(gems); i++ {
		for j := i + 1; j < len(gems); j++ {
			for k := j + 1; k < len(gems); k++ {
				count := items.GemMap{
					gems[i]: 1,
					gems[j]: 1,
					gems[k]: 1,
				}
				sets = append(sets, count.ToCount())
			}
		}
	}
	return sets
}
