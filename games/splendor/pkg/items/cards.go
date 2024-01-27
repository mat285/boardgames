package items

func Cards() []Card {
	cards := append(
		LevelOneCards(),
		append(LevelTwoCards(),
			LevelThreeCards()...,
		)...,
	)
	for i := range cards {
		cards[i].ID = i
	}
	return cards
}

func LevelOneCards() []Card {
	return []Card{
		// ** Level 1**
		// Black
		{
			Level: 0,
			Value: 0,
			Type:  GemObsidian,
			Cost: GemCount{
				Diamond:  1,
				Sapphire: 1,
				Emerald:  1,
				Ruby:     1,
				Obsidian: 0,
			},
		},
		{
			Level: 0,
			Value: 0,
			Type:  GemObsidian,
			Cost: GemCount{
				Diamond:  1,
				Sapphire: 2,
				Emerald:  1,
				Ruby:     1,
				Obsidian: 0,
			},
		},
		{
			Level: 0,
			Value: 0,
			Type:  GemObsidian,
			Cost: GemCount{
				Diamond:  2,
				Sapphire: 2,
				Emerald:  0,
				Ruby:     1,
				Obsidian: 0,
			},
		},
		{
			Level: 0,
			Value: 0,
			Type:  GemObsidian,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 0,
				Emerald:  1,
				Ruby:     3,
				Obsidian: 1,
			},
		},
		{
			Level: 0,
			Value: 0,
			Type:  GemObsidian,
			Cost: GemCount{
				Diamond:  2,
				Sapphire: 0,
				Emerald:  2,
				Ruby:     0,
				Obsidian: 0,
			},
		},
		{
			Level: 0,
			Value: 0,
			Type:  GemObsidian,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 0,
				Emerald:  3,
				Ruby:     0,
				Obsidian: 0,
			},
		},
		{
			Level: 0,
			Value: 1,
			Type:  GemObsidian,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 4,
				Emerald:  0,
				Ruby:     0,
				Obsidian: 0,
			},
		},
		// Blue
		{
			Level: 0,
			Value: 0,
			Type:  GemSapphire,
			Cost: GemCount{
				Diamond:  1,
				Sapphire: 0,
				Emerald:  1,
				Ruby:     1,
				Obsidian: 1,
			},
		},
		{
			Level: 0,
			Value: 0,
			Type:  GemSapphire,
			Cost: GemCount{
				Diamond:  1,
				Sapphire: 0,
				Emerald:  1,
				Ruby:     2,
				Obsidian: 1,
			},
		},
		{
			Level: 0,
			Value: 0,
			Type:  GemSapphire,
			Cost: GemCount{
				Diamond:  1,
				Sapphire: 0,
				Emerald:  2,
				Ruby:     2,
				Obsidian: 0,
			},
		},
		{
			Level: 0,
			Value: 0,
			Type:  GemSapphire,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 1,
				Emerald:  3,
				Ruby:     1,
				Obsidian: 0,
			},
		},
		{
			Level: 0,
			Value: 0,
			Type:  GemSapphire,
			Cost: GemCount{
				Diamond:  1,
				Sapphire: 0,
				Emerald:  0,
				Ruby:     0,
				Obsidian: 2,
			},
		},
		{
			Level: 0,
			Value: 0,
			Type:  GemSapphire,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 0,
				Emerald:  2,
				Ruby:     0,
				Obsidian: 2,
			},
		},
		{
			Level: 0,
			Value: 1,
			Type:  GemSapphire,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 0,
				Emerald:  0,
				Ruby:     4,
				Obsidian: 0,
			},
		},
		// white
		{
			Level: 0,
			Value: 0,
			Type:  GemDiamond,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 1,
				Emerald:  1,
				Ruby:     1,
				Obsidian: 1,
			},
		},
		{
			Level: 0,
			Value: 0,
			Type:  GemDiamond,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 1,
				Emerald:  2,
				Ruby:     1,
				Obsidian: 1,
			},
		},
		{
			Level: 0,
			Value: 0,
			Type:  GemDiamond,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 2,
				Emerald:  2,
				Ruby:     0,
				Obsidian: 1,
			},
		},
		{
			Level: 0,
			Value: 0,
			Type:  GemDiamond,
			Cost: GemCount{
				Diamond:  3,
				Sapphire: 1,
				Emerald:  0,
				Ruby:     0,
				Obsidian: 1,
			},
		},
		{
			Level: 0,
			Value: 0,
			Type:  GemDiamond,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 0,
				Emerald:  2,
				Ruby:     0,
				Obsidian: 1,
			},
		},
		{
			Level: 0,
			Value: 0,
			Type:  GemDiamond,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 2,
				Emerald:  0,
				Ruby:     0,
				Obsidian: 2,
			},
		},
		{
			Level: 0,
			Value: 0,
			Type:  GemDiamond,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 3,
				Emerald:  0,
				Ruby:     0,
				Obsidian: 0,
			},
		},
		{
			Level: 0,
			Value: 1,
			Type:  GemDiamond,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 0,
				Emerald:  4,
				Ruby:     0,
				Obsidian: 0,
			},
		},
		// green
		{
			Level: 0,
			Value: 0,
			Type:  GemEmerald,
			Cost: GemCount{
				Diamond:  1,
				Sapphire: 1,
				Emerald:  0,
				Ruby:     1,
				Obsidian: 1,
			},
		},
		{
			Level: 0,
			Value: 0,
			Type:  GemEmerald,
			Cost: GemCount{
				Diamond:  1,
				Sapphire: 1,
				Emerald:  0,
				Ruby:     1,
				Obsidian: 2,
			},
		},
		{
			Level: 0,
			Value: 0,
			Type:  GemEmerald,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 1,
				Emerald:  0,
				Ruby:     2,
				Obsidian: 2,
			},
		},
		{
			Level: 0,
			Value: 0,
			Type:  GemEmerald,
			Cost: GemCount{
				Diamond:  1,
				Sapphire: 3,
				Emerald:  1,
				Ruby:     0,
				Obsidian: 0,
			},
		},
		{
			Level: 0,
			Value: 0,
			Type:  GemEmerald,
			Cost: GemCount{
				Diamond:  2,
				Sapphire: 1,
				Emerald:  0,
				Ruby:     0,
				Obsidian: 0,
			},
		},
		{
			Level: 0,
			Value: 0,
			Type:  GemEmerald,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 2,
				Emerald:  0,
				Ruby:     2,
				Obsidian: 0,
			},
		},
		{
			Level: 0,
			Value: 0,
			Type:  GemEmerald,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 0,
				Emerald:  0,
				Ruby:     3,
				Obsidian: 0,
			},
		},
		{
			Level: 0,
			Value: 1,
			Type:  GemEmerald,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 0,
				Emerald:  0,
				Ruby:     0,
				Obsidian: 4,
			},
		},
		// red
		{
			Level: 0,
			Value: 0,
			Type:  GemRuby,
			Cost: GemCount{
				Diamond:  1,
				Sapphire: 1,
				Emerald:  1,
				Ruby:     0,
				Obsidian: 1,
			},
		},
		{
			Level: 0,
			Value: 0,
			Type:  GemRuby,
			Cost: GemCount{
				Diamond:  2,
				Sapphire: 1,
				Emerald:  1,
				Ruby:     0,
				Obsidian: 1,
			},
		},
		{
			Level: 0,
			Value: 0,
			Type:  GemRuby,
			Cost: GemCount{
				Diamond:  2,
				Sapphire: 0,
				Emerald:  1,
				Ruby:     0,
				Obsidian: 2,
			},
		},
		{
			Level: 0,
			Value: 0,
			Type:  GemRuby,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 2,
				Emerald:  1,
				Ruby:     0,
				Obsidian: 0,
			},
		},
		{
			Level: 0,
			Value: 0,
			Type:  GemRuby,
			Cost: GemCount{
				Diamond:  2,
				Sapphire: 0,
				Emerald:  0,
				Ruby:     2,
				Obsidian: 0,
			},
		},
		{
			Level: 0,
			Value: 0,
			Type:  GemRuby,
			Cost: GemCount{
				Diamond:  3,
				Sapphire: 0,
				Emerald:  0,
				Ruby:     0,
				Obsidian: 0,
			},
		},
		{
			Level: 0,
			Value: 1,
			Type:  GemRuby,
			Cost: GemCount{
				Diamond:  4,
				Sapphire: 0,
				Emerald:  0,
				Ruby:     0,
				Obsidian: 0,
			},
		},
	}
}

func LevelTwoCards() []Card {
	return []Card{
		// **Level 2**
		// Black
		{
			Level: 1,
			Value: 1,
			Type:  GemObsidian,
			Cost: GemCount{
				Diamond:  3,
				Sapphire: 2,
				Emerald:  2,
				Ruby:     0,
				Obsidian: 0,
			},
		},
		{
			Level: 1,
			Value: 1,
			Type:  GemObsidian,
			Cost: GemCount{
				Diamond:  3,
				Sapphire: 0,
				Emerald:  3,
				Ruby:     0,
				Obsidian: 2,
			},
		},
		{
			Level: 1,
			Value: 2,
			Type:  GemObsidian,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 1,
				Emerald:  4,
				Ruby:     2,
				Obsidian: 0,
			},
		},
		{
			Level: 1,
			Value: 2,
			Type:  GemObsidian,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 0,
				Emerald:  5,
				Ruby:     3,
				Obsidian: 0,
			},
		},
		{
			Level: 1,
			Value: 2,
			Type:  GemObsidian,
			Cost: GemCount{
				Diamond:  5,
				Sapphire: 0,
				Emerald:  0,
				Ruby:     0,
				Obsidian: 0,
			},
		},
		{
			Level: 1,
			Value: 3,
			Type:  GemObsidian,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 0,
				Emerald:  0,
				Ruby:     0,
				Obsidian: 6,
			},
		},

		// blue
		{
			Level: 1,
			Value: 1,
			Type:  GemSapphire,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 2,
				Emerald:  2,
				Ruby:     3,
				Obsidian: 0,
			},
		},
		{
			Level: 1,
			Value: 1,
			Type:  GemSapphire,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 2,
				Emerald:  3,
				Ruby:     0,
				Obsidian: 3,
			},
		},
		{
			Level: 1,
			Value: 2,
			Type:  GemSapphire,
			Cost: GemCount{
				Diamond:  5,
				Sapphire: 3,
				Emerald:  0,
				Ruby:     0,
				Obsidian: 0,
			},
		},
		{
			Level: 1,
			Value: 2,
			Type:  GemSapphire,
			Cost: GemCount{
				Diamond:  2,
				Sapphire: 0,
				Emerald:  0,
				Ruby:     1,
				Obsidian: 4,
			},
		},
		{
			Level: 1,
			Value: 2,
			Type:  GemSapphire,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 5,
				Emerald:  0,
				Ruby:     0,
				Obsidian: 0,
			},
		},
		{
			Level: 1,
			Value: 3,
			Type:  GemSapphire,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 6,
				Emerald:  0,
				Ruby:     0,
				Obsidian: 0,
			},
		},

		// white
		{
			Level: 1,
			Value: 1,
			Type:  GemDiamond,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 0,
				Emerald:  3,
				Ruby:     2,
				Obsidian: 2,
			},
		},
		{
			Level: 1,
			Value: 1,
			Type:  GemDiamond,
			Cost: GemCount{
				Diamond:  2,
				Sapphire: 3,
				Emerald:  0,
				Ruby:     3,
				Obsidian: 0,
			},
		},
		{
			Level: 1,
			Value: 2,
			Type:  GemDiamond,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 0,
				Emerald:  1,
				Ruby:     4,
				Obsidian: 2,
			},
		},
		{
			Level: 1,
			Value: 2,
			Type:  GemDiamond,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 0,
				Emerald:  0,
				Ruby:     5,
				Obsidian: 3,
			},
		},
		{
			Level: 1,
			Value: 2,
			Type:  GemDiamond,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 0,
				Emerald:  0,
				Ruby:     5,
				Obsidian: 0,
			},
		},
		{
			Level: 1,
			Value: 3,
			Type:  GemDiamond,
			Cost: GemCount{
				Diamond:  6,
				Sapphire: 0,
				Emerald:  0,
				Ruby:     0,
				Obsidian: 0,
			},
		},

		// green
		{
			Level: 1,
			Value: 1,
			Type:  GemEmerald,
			Cost: GemCount{
				Diamond:  3,
				Sapphire: 0,
				Emerald:  2,
				Ruby:     3,
				Obsidian: 0,
			},
		},
		{
			Level: 1,
			Value: 1,
			Type:  GemEmerald,
			Cost: GemCount{
				Diamond:  2,
				Sapphire: 3,
				Emerald:  0,
				Ruby:     0,
				Obsidian: 2,
			},
		},
		{
			Level: 1,
			Value: 2,
			Type:  GemEmerald,
			Cost: GemCount{
				Diamond:  4,
				Sapphire: 2,
				Emerald:  0,
				Ruby:     0,
				Obsidian: 1,
			},
		},
		{
			Level: 1,
			Value: 2,
			Type:  GemEmerald,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 5,
				Emerald:  3,
				Ruby:     0,
				Obsidian: 0,
			},
		},
		{
			Level: 1,
			Value: 2,
			Type:  GemEmerald,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 0,
				Emerald:  5,
				Ruby:     0,
				Obsidian: 0,
			},
		},
		{
			Level: 1,
			Value: 3,
			Type:  GemEmerald,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 0,
				Emerald:  6,
				Ruby:     0,
				Obsidian: 0,
			},
		},

		// red
		{
			Level: 1,
			Value: 1,
			Type:  GemRuby,
			Cost: GemCount{
				Diamond:  2,
				Sapphire: 0,
				Emerald:  0,
				Ruby:     2,
				Obsidian: 3,
			},
		},
		{
			Level: 1,
			Value: 1,
			Type:  GemRuby,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 3,
				Emerald:  0,
				Ruby:     2,
				Obsidian: 3,
			},
		},
		{
			Level: 1,
			Value: 2,
			Type:  GemRuby,
			Cost: GemCount{
				Diamond:  1,
				Sapphire: 4,
				Emerald:  2,
				Ruby:     0,
				Obsidian: 0,
			},
		},
		{
			Level: 1,
			Value: 2,
			Type:  GemRuby,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 0,
				Emerald:  0,
				Ruby:     0,
				Obsidian: 5,
			},
		},
		{
			Level: 1,
			Value: 3,
			Type:  GemRuby,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 0,
				Emerald:  0,
				Ruby:     6,
				Obsidian: 0,
			},
		},
	}
}

func LevelThreeCards() []Card {
	return []Card{
		// **Level 3**
		// Black
		{
			Level: 2,
			Value: 3,
			Type:  GemObsidian,
			Cost: GemCount{
				Diamond:  3,
				Sapphire: 3,
				Emerald:  5,
				Ruby:     3,
				Obsidian: 0,
			},
		},
		{
			Level: 2,
			Value: 4,
			Type:  GemObsidian,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 0,
				Emerald:  0,
				Ruby:     7,
				Obsidian: 0,
			},
		},
		{
			Level: 2,
			Value: 4,
			Type:  GemObsidian,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 0,
				Emerald:  3,
				Ruby:     6,
				Obsidian: 3,
			},
		},
		{
			Level: 2,
			Value: 5,
			Type:  GemObsidian,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 0,
				Emerald:  0,
				Ruby:     7,
				Obsidian: 3,
			},
		},

		// blue
		{
			Level: 2,
			Value: 3,
			Type:  GemSapphire,
			Cost: GemCount{
				Diamond:  3,
				Sapphire: 0,
				Emerald:  3,
				Ruby:     3,
				Obsidian: 5,
			},
		},
		{
			Level: 2,
			Value: 4,
			Type:  GemSapphire,
			Cost: GemCount{
				Diamond:  7,
				Sapphire: 0,
				Emerald:  0,
				Ruby:     0,
				Obsidian: 0,
			},
		},
		{
			Level: 2,
			Value: 4,
			Type:  GemSapphire,
			Cost: GemCount{
				Diamond:  6,
				Sapphire: 3,
				Emerald:  0,
				Ruby:     0,
				Obsidian: 3,
			},
		},
		{
			Level: 2,
			Value: 5,
			Type:  GemSapphire,
			Cost: GemCount{
				Diamond:  7,
				Sapphire: 3,
				Emerald:  0,
				Ruby:     0,
				Obsidian: 0,
			},
		},

		// white
		{
			Level: 2,
			Value: 3,
			Type:  GemDiamond,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 3,
				Emerald:  3,
				Ruby:     5,
				Obsidian: 3,
			},
		},
		{
			Level: 2,
			Value: 4,
			Type:  GemDiamond,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 0,
				Emerald:  0,
				Ruby:     0,
				Obsidian: 7,
			},
		},
		{
			Level: 2,
			Value: 4,
			Type:  GemDiamond,
			Cost: GemCount{
				Diamond:  3,
				Sapphire: 0,
				Emerald:  0,
				Ruby:     0,
				Obsidian: 7,
			},
		},
		{
			Level: 2,
			Value: 5,
			Type:  GemDiamond,
			Cost: GemCount{
				Diamond:  3,
				Sapphire: 0,
				Emerald:  0,
				Ruby:     0,
				Obsidian: 7,
			},
		},

		// green
		{
			Level: 2,
			Value: 3,
			Type:  GemEmerald,
			Cost: GemCount{
				Diamond:  5,
				Sapphire: 3,
				Emerald:  0,
				Ruby:     3,
				Obsidian: 3,
			},
		},
		{
			Level: 2,
			Value: 4,
			Type:  GemEmerald,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 7,
				Emerald:  0,
				Ruby:     0,
				Obsidian: 0,
			},
		},
		{
			Level: 2,
			Value: 4,
			Type:  GemEmerald,
			Cost: GemCount{
				Diamond:  3,
				Sapphire: 6,
				Emerald:  3,
				Ruby:     0,
				Obsidian: 0,
			},
		},
		{
			Level: 2,
			Value: 5,
			Type:  GemEmerald,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 7,
				Emerald:  3,
				Ruby:     0,
				Obsidian: 0,
			},
		},

		// red
		{
			Level: 2,
			Value: 3,
			Type:  GemRuby,
			Cost: GemCount{
				Diamond:  3,
				Sapphire: 5,
				Emerald:  3,
				Ruby:     0,
				Obsidian: 3,
			},
		},
		{
			Level: 2,
			Value: 4,
			Type:  GemRuby,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 7,
				Emerald:  0,
				Ruby:     0,
				Obsidian: 0,
			},
		},
		{
			Level: 2,
			Value: 4,
			Type:  GemRuby,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 3,
				Emerald:  6,
				Ruby:     3,
				Obsidian: 0,
			},
		},
		{
			Level: 2,
			Value: 5,
			Type:  GemRuby,
			Cost: GemCount{
				Diamond:  0,
				Sapphire: 0,
				Emerald:  7,
				Ruby:     3,
				Obsidian: 0,
			},
		},
	}
}
