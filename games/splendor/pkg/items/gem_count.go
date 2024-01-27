package items

type GemMap map[Gem]int

type GemCount struct {
	Diamond  int
	Sapphire int
	Emerald  int
	Ruby     int
	Obsidian int
	Wild     int
}

func GemCountFromMap(counts GemMap) GemCount {
	return GemCount{
		Diamond:  counts[GemDiamond],
		Sapphire: counts[GemSapphire],
		Emerald:  counts[GemEmerald],
		Ruby:     counts[GemRuby],
		Obsidian: counts[GemObsidian],
		Wild:     counts[GemWild],
	}
}

func (gc GemCount) Total() int {
	return gc.Diamond + gc.Sapphire + gc.Emerald + gc.Ruby + gc.Obsidian + gc.Wild
}

func (gc GemCount) ToMap() GemMap {
	return GemMap{
		GemDiamond:  gc.Diamond,
		GemSapphire: gc.Sapphire,
		GemEmerald:  gc.Emerald,
		GemRuby:     gc.Ruby,
		GemObsidian: gc.Obsidian,
		GemWild:     gc.Wild,
	}
}

func (gm GemMap) ToCount() GemCount {
	return GemCountFromMap(gm)
}

func (gc GemCount) Get(gem Gem) int {
	var count int
	gc.apply(gem, func(_ Gem, curr int) int {
		count = curr
		return curr
	})
	return count
}

func (gc GemCount) AddGem(gem Gem, count int) GemCount {
	return gc.apply(gem, func(_ Gem, curr int) int {
		return curr + count
	})
}

func (gc GemCount) apply(gem Gem, fn func(Gem, int) int) GemCount {
	switch gem {
	case GemDiamond:
		gc.Diamond = fn(GemDiamond, gc.Diamond)
	case GemSapphire:
		gc.Sapphire = fn(GemSapphire, gc.Sapphire)
	case GemEmerald:
		gc.Emerald = fn(GemEmerald, gc.Emerald)
	case GemRuby:
		gc.Ruby = fn(GemRuby, gc.Ruby)
	case GemObsidian:
		gc.Obsidian = fn(GemObsidian, gc.Obsidian)
	case GemWild:
		gc.Wild = fn(GemWild, gc.Wild)
	}
	return gc
}

func (gc GemCount) Add(other GemCount) GemCount {
	gc.Diamond += other.Diamond
	gc.Sapphire += other.Sapphire
	gc.Emerald += other.Emerald
	gc.Ruby += other.Ruby
	gc.Obsidian += other.Obsidian
	gc.Wild += other.Wild
	return gc
}

func (gc GemCount) Sub(other GemCount) GemCount {
	gc.Diamond -= other.Diamond
	gc.Sapphire -= other.Sapphire
	gc.Emerald -= other.Emerald
	gc.Ruby -= other.Ruby
	gc.Obsidian -= other.Obsidian
	gc.Wild -= other.Wild
	return gc
}

func (gc GemCount) Equal(other GemCount) bool {
	return gc.Diamond == other.Diamond &&
		gc.Sapphire == other.Sapphire &&
		gc.Emerald == other.Emerald &&
		gc.Ruby == other.Ruby &&
		gc.Obsidian == other.Obsidian &&
		gc.Wild == other.Wild
}

func (gc GemCount) LessThanEqual(other GemCount) bool {
	return gc.Diamond <= other.Diamond &&
		gc.Sapphire <= other.Sapphire &&
		gc.Emerald <= other.Emerald &&
		gc.Ruby <= other.Ruby &&
		gc.Obsidian <= other.Obsidian &&
		gc.Wild <= other.Wild
}

func (gm GemMap) Keys() []Gem {
	keys := make([]Gem, 0, len(gm))
	for k := range gm {
		keys = append(keys, k)
	}
	return keys
}

func (gm GemMap) NonNegative() bool {
	for _, c := range gm {
		if c < 0 {
			return false
		}
	}
	return true
}

func (gm GemMap) Total() int {
	count := 0
	for _, c := range gm {
		count += c
	}
	return count
}

func (gm GemMap) Sub(other GemMap) {
	for k, v := range other {
		gm[k] -= v
	}
}
