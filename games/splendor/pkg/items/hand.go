package items

type Hand struct {
	Gems  GemCount
	Cards []Card
	Bonus []Bonus

	Reserved []Card
}

func NewHand() Hand {
	return Hand{
		Gems:     GemCount{},
		Cards:    make([]Card, 0),
		Bonus:    make([]Bonus, 0),
		Reserved: make([]Card, 0),
	}
}

func (h Hand) Points() int {
	points := 0

	for _, card := range h.Cards {
		points += card.Value
	}
	for _, bonus := range h.Bonus {
		points += bonus.Value
	}

	return points
}

func (h Hand) Collect(gems GemCount) Hand {
	h.Gems = h.Gems.Add(gems)
	return h
}

func (h Hand) BonusesEarned(current []Bonus) []Bonus {
	bs := []Bonus{}
	cc := h.CardCounts()
	for _, b := range current {
		if b.Count.LessThanEqual(cc) {
			bs = append(bs, b)
		}
	}
	return bs
}

func (h Hand) CardCounts() GemCount {
	gc := GemCount{}
	for _, card := range h.Cards {
		gc = gc.AddGem(card.Type, 1)
	}
	return gc
}

func (h Hand) CanReserve() bool {
	return len(h.Reserved) < 3
}

func (h Hand) Reserve(card Card) Hand {
	h.Gems = h.Gems.AddGem(GemWild, 1)
	h.Reserved = append(h.Reserved, card)
	return h
}

func (h Hand) CanPurchase(card Card) bool {
	price := card.Price().ToMap()
	cardCounts := h.CardCounts().ToMap()

	for g := range price {
		price[g] -= cardCounts[g]
		if price[g] <= 0 {
			delete(price, g)
		}
	}

	for gem, count := range h.Gems.ToMap() {
		if price[gem] <= count {
			delete(price, gem)
		}
	}

	remaining := 0
	for _, count := range price {
		remaining += count
	}

	return remaining <= h.Gems.Wild
}

func (h Hand) reservedWithout(card Card) []Card {
	reserved := make([]Card, 0, len(h.Reserved))
	for _, r := range h.Reserved {
		if card.Equal(r) {
			continue
		}
		reserved = append(reserved, r)
	}
	return reserved
}

func (h Hand) Purchase(card Card) Hand {
	price := card.Price().ToMap()
	cardCounts := h.CardCounts().ToMap()

	for g := range price {
		price[g] -= cardCounts[g]
		if price[g] <= 0 {
			delete(price, g)
		}
	}

	gems := h.Gems.ToMap()

	for g, c := range price {
		if gems[g] >= c {
			gems[g] -= c
		} else {
			gems[GemWild] -= (c - gems[g])
			gems[g] = 0
		}
	}

	h.Cards = append(h.Cards, card)

	return Hand{
		Gems:     GemCountFromMap(gems),
		Cards:    h.Cards,
		Bonus:    h.Bonus,
		Reserved: h.reservedWithout(card),
	}
}
