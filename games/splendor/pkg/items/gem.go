package items

import (
	"fmt"
	"strings"
)

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

func ParseGem(str string) (Gem, error) {
	switch strings.ToLower(str) {
	case string(GemDiamond), "white":
		return GemDiamond, nil
	case string(GemSapphire), "blue":
		return GemSapphire, nil
	case string(GemEmerald), "green":
		return GemEmerald, nil
	case string(GemRuby), "red":
		return GemRuby, nil
	case string(GemObsidian), "black":
		return GemObsidian, nil
	default:
		return "", fmt.Errorf("Not a gem")
	}
}
