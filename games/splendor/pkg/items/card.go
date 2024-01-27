package items

type Card struct {
	ID    int
	Level int
	Value int
	Type  Gem
	Cost  GemCount
}

func (c Card) Price() GemCount {
	return c.Cost
}

func (c Card) Equal(other Card) bool {
	if c.Value != other.Value || c.Type != other.Type {
		return false
	}
	return c.Cost.Equal(other.Cost)
}

func ContainsCard(cards []Card, card Card) bool {
	for _, c := range cards {
		if c.Equal(card) {
			return true
		}
	}
	return false
}

func CloneCards(cards []Card) []Card {
	ret := make([]Card, len(cards))
	for i := range cards {
		ret[i] = cards[i]
	}
	return ret
}
