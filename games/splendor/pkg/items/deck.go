package items

import "math/rand"

type Deck struct {
	Shown []Card
	Pile  []Card
}

func NewDeck(cards []Card) Deck {
	return Deck{
		Shown: make([]Card, 0),
		Pile:  cards,
	}
}

func (d Deck) Deal(num int) Deck {
	pile := d.Pile
	shown := CloneCards(d.Shown)

	for i := 0; i < num; i++ {
		var card Card
		card, pile = PickRandomCard(pile)
		shown = append(shown, card)
	}

	d.Shown = shown
	d.Pile = pile
	return d
}

func (d Deck) RemoveAndReplace(card Card) Deck {
	shown := make([]Card, 0, len(d.Shown))
	for _, s := range d.Shown {
		if s.Equal(card) {
			continue
		}
		shown = append(shown, s)
	}

	d.Shown = shown
	return d.Deal(1)
}

func PickRandomCard(cards []Card) (Card, []Card) {
	if len(cards) == 0 {
		return Card{}, []Card{}
	}
	i := rand.Intn(len(cards))
	card := cards[i]
	last := len(cards) - 1

	cards[i] = cards[last]
	return card, cards[:last]
}
