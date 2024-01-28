package items

import "math/rand"

type Bonus struct {
	Value int
	Count GemCount
}

func ContainsBonus(bs []Bonus, bonus Bonus) bool {
	for _, b := range bs {
		if b == bonus {
			return true
		}
	}
	return false
}

func CloneBonuses(bs []Bonus) []Bonus {
	ret := make([]Bonus, len(bs))
	for i := range bs {
		ret[i] = bs[i]
	}
	return ret
}

func RandomBonuses(n int) []Bonus {
	all := Bonuses()
	if n <= 0 {
		return []Bonus{}
	}
	if n >= len(all) {
		return all
	}
	ret := make([]Bonus, n)
	for i := 0; i < n; i++ {
		ret[i], all = RandomBonus(all)
	}
	return ret
}

func RandomBonus(bs []Bonus) (Bonus, []Bonus) {
	if len(bs) == 0 {
		return Bonus{}, []Bonus{}
	}
	i := rand.Intn(len(bs))
	bonus := bs[i]
	last := len(bs) - 1

	bs[i] = bs[last]
	return bonus, bs[:last]
}

func Bonuses() []Bonus {
	return []Bonus{
		{
			Value: 3,
			Count: GemCount{
				Diamond:  3,
				Sapphire: 3,
				Emerald:  3,
				Ruby:     0,
				Obsidian: 0,
			},
		},
		{
			Value: 3,
			Count: GemCount{
				Diamond:  0,
				Sapphire: 3,
				Emerald:  3,
				Ruby:     3,
				Obsidian: 0,
			},
		},
		{
			Value: 3,
			Count: GemCount{
				Diamond:  3,
				Sapphire: 0,
				Emerald:  0,
				Ruby:     3,
				Obsidian: 3,
			},
		},
		{
			Value: 3,
			Count: GemCount{
				Diamond:  3,
				Sapphire: 3,
				Emerald:  0,
				Ruby:     0,
				Obsidian: 3,
			},
		},
		{
			Value: 3,
			Count: GemCount{
				Diamond:  0,
				Sapphire: 0,
				Emerald:  3,
				Ruby:     3,
				Obsidian: 3,
			},
		},
		{
			Value: 3,
			Count: GemCount{
				Diamond:  0,
				Sapphire: 0,
				Emerald:  0,
				Ruby:     4,
				Obsidian: 4,
			},
		},
		{
			Value: 3,
			Count: GemCount{
				Diamond:  4,
				Sapphire: 0,
				Emerald:  0,
				Ruby:     0,
				Obsidian: 4,
			},
		},
		{
			Value: 3,
			Count: GemCount{
				Diamond:  0,
				Sapphire: 0,
				Emerald:  4,
				Ruby:     4,
				Obsidian: 0,
			},
		},
		{
			Value: 3,
			Count: GemCount{
				Diamond:  0,
				Sapphire: 4,
				Emerald:  4,
				Ruby:     0,
				Obsidian: 0,
			},
		},
		{
			Value: 3,
			Count: GemCount{
				Diamond:  4,
				Sapphire: 4,
				Emerald:  0,
				Ruby:     0,
				Obsidian: 0,
			},
		},
	}
}
