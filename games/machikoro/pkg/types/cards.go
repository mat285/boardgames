package types

import common "github.com/mat285/boardgames/pkg/common/v1alpha1"

type CardCount struct {
	Card
	Count int
}

func Cards() []Card {
	counts := CardCounts()
	id := 0
	cards := make([]Card, 0, len(counts)*5)
	for _, count := range counts {
		for i := 0; i <= count.Count; i++ {
			card := count.Card
			card.ID = id
			cards = append(cards, card)
			id++
		}
	}
	return cards
}

func sortPiles(cards []Card) [][]Card {
	ret := make([][]Card, 3)
	for _, card := range cards {
		index := -1
		if card.Type == CardTypeLandMark {
			index = 2
		} else if card.Activation.Low < 7 {
			index = 0
		} else {
			index = 1
		}
		ret[index] = append(ret[index], card)
	}
	return ret
}

func CardCounts() []CardCount {
	return []CardCount{
		{
			Card: Card{
				Name:        "Wheat Field",
				Type:        CardTypePrimaryIndustry,
				Icon:        CardIconWheat,
				Cost:        1,
				Activation:  common.NewIntRange(1, 2),
				Description: "",
			},
			Count: 6,
		},
	}
}
