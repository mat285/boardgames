package items

type Gem string

const (
	GemDiamond  Gem = "diamond"
	GemObsidian Gem = "obsidian"
	GemRuby     Gem = "ruby"
	GemSapphire Gem = "sapphire"
	GemEmerald  Gem = "emerald"

	GemWild Gem = "wild"

	GemWhite = GemDiamond
	GemBlack = GemObsidian
	GemRed   = GemRuby
	GemBlue  = GemSapphire
	GemGreen = GemEmerald
)

func Gems() GemCount {
	return GemCount{
		Diamond:  7,
		Sapphire: 7,
		Emerald:  7,
		Ruby:     7,
		Obsidian: 7,
		Wild:     5,
	}
}
