package items_test

import (
	"testing"

	"github.com/blend/go-sdk/assert"
	"github.com/mat285/boardgames/games/splendor/pkg/items"
)

func TestCanPurchase(t *testing.T) {
	it := assert.New(t)

	h := items.Hand{
		Gems: items.GemCountFromMap(
			items.GemMap{
				items.GemBlack: 2,
				items.GemRed:   2,
				items.GemWild:  1,
			},
		),
	}

	c := items.Card{
		Cost: items.GemCountFromMap(
			items.GemMap{
				items.GemBlack: 2,
				items.GemRed:   3,
			},
		),
	}

	it.True(h.CanPurchase(c))
}

func TestPurchase(t *testing.T) {
	it := assert.New(t)

	h := items.Hand{
		Gems: items.GemCountFromMap(
			items.GemMap{
				items.GemBlack: 2,
				items.GemRed:   2,
				items.GemWild:  1,
			},
		),
	}

	c := items.Card{
		Cost: items.GemCountFromMap(
			items.GemMap{
				items.GemBlack: 2,
				items.GemRed:   3,
			},
		),
	}

	nh := h.Purchase(c)
	it.Len(nh.Cards, 1)
	it.Equal(0, nh.Gems.Total())
}
