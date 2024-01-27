package items

type Board struct {
	Gems       GemCount
	LevelOne   Deck
	LevelTwo   Deck
	LevelThree Deck
	Bonuses    []Bonus
}

func NewBoard() Board {
	return Board{
		Gems:       Gems(),
		LevelOne:   NewDeck(LevelOneCards()).Deal(4),
		LevelTwo:   NewDeck(LevelTwoCards()).Deal(4),
		LevelThree: NewDeck(LevelThreeCards()).Deal(4),
		Bonuses:    Bonuses(),
	}

}

func (b Board) IsCardOnBoard(card Card) bool {
	switch card.Level {
	case 0:
		return ContainsCard(b.LevelOne.Shown, card)
	case 1:
		return ContainsCard(b.LevelTwo.Shown, card)
	case 2:
		return ContainsCard(b.LevelThree.Shown, card)
	default:
		return false
	}
}

func (b Board) AvailableCards() []Card {
	ret := make([]Card, 0, len(b.LevelOne.Shown)+len(b.LevelTwo.Shown)+len(b.LevelThree.Shown))
	ret = append(ret, b.LevelOne.Shown...)
	ret = append(ret, b.LevelTwo.Shown...)
	ret = append(ret, b.LevelThree.Shown...)
	return ret
}

func (b Board) RemoveCard(card Card) Board {
	switch card.Level {
	case 0:
		b.LevelOne = b.LevelOne.RemoveAndReplace(card)
	case 1:
		b.LevelTwo = b.LevelTwo.RemoveAndReplace(card)
	case 2:
		b.LevelThree = b.LevelThree.RemoveAndReplace(card)
	default:
	}
	return b
}

func (b Board) RemoveBonuses(bs []Bonus) Board {
	nb := []Bonus{}
	for _, curr := range b.Bonuses {
		if !ContainsBonus(bs, curr) {
			nb = append(nb, curr)
		}
	}
	b.Bonuses = nb
	return b
}
