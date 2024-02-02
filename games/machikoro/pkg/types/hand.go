package types

type Hand struct {
	Money int
	Cards map[CardType][]Card
}

func NewHand(start int) Hand {
	return Hand{
		Money: start,
		Cards: make(map[CardType][]Card),
	}
}

func (h Hand) LandmarksCount() int {
	return len(h.Cards[CardTypeLandMark])
}

func (h Hand) Collect(money int) Hand {
	h.Money += money
	h.Cards = cloneHandCards(h.Cards)
	return h
}

func (h Hand) Subtract(money int) Hand {
	h.Money -= money
	h.Cards = cloneHandCards(h.Cards)
	return h
}

func (h Hand) Purchase(card Card) Hand {
	h.Money -= card.GetCost(LandmarkModifier(h.LandmarksCount()))
	h.Cards = cloneHandCards(h.Cards)
	h.Cards[card.Type] = append(h.Cards[card.Type], card)
	return h
}

func cloneHandCards(cards map[CardType][]Card) map[CardType][]Card {
	ret := make(map[CardType][]Card)
	for k, v := range cards {
		for _, c := range v {
			ret[k] = append(ret[k], c)
		}
	}
	return ret
}
