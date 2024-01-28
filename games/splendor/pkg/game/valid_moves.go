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

	gems := s.Board.Gems.ToMap()
	delete(gems, items.GemWild)
	sets := allSetsOf3(gems.Keys())

	for _, set := range sets {
		moves = append(moves, &Move{Collect: &CollectMove{Take: set}})
	}

	for k, v := range gems {
		if v >= items.MinGemsToTakeTwo {
			count := items.GemMap{k: 2}
			moves = append(moves, &Move{Collect: &CollectMove{Take: count.ToCount()}})
		}
	}
	takeret := make([]*Move, 0, len(moves))
	hand := s.Players[s.CurrentIndex].Hand
	total := hand.Gems.Total()
	gemSlice := hand.Gems.ToSlice()
	for i, move := range moves {
		count := move.Collect.Take.Total()
		if count+total > items.MaxGems {
			n := count + total - items.MaxGems
			sets := allSetsOfN(n, gemSlice)
			for _, set := range sets {
				takeret = append(takeret, &Move{
					Collect: &CollectMove{
						Take:   move.Collect.Take,
						Return: set,
					},
				})
			}
		} else {
			takeret = append(takeret, moves[i])
		}
	}
	return takeret
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

func allSetsOfN(n int, gems []items.Gem) []items.GemCount {
	switch n {
	case 1:
		return allSetsOf1(gems)
	case 2:
		return allSetsOf2(gems)
	case 3:
		return allSetsOf3(gems)
	}
	return nil
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

func allSetsOf2(gems []items.Gem) []items.GemCount {
	sets := make([]items.GemCount, 0)
	for i := 0; i < len(gems); i++ {
		for j := i + 1; j < len(gems); j++ {
			count := items.GemMap{
				gems[i]: 1,
				gems[j]: 1,
			}
			sets = append(sets, count.ToCount())
		}
	}
	return sets
}

func allSetsOf1(gems []items.Gem) []items.GemCount {

	m := make(map[items.Gem]items.GemCount)
	for _, gem := range gems {
		c := items.GemCount{}
		c.AddGem(gem, 1)
		m[gem] = c
	}
	sets := make([]items.GemCount, 0, len(m))

	for _, gc := range m {
		sets = append(sets, gc)
	}
	return sets

}
