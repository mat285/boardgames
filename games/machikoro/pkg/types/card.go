package types

import (
	common "github.com/mat285/boardgames/pkg/common/v1alpha1"
)

type Card struct {
	ID          int
	Name        string
	Type        CardType
	Icon        CardIcon
	Cost        Cost
	Activation  common.IntRange
	Description string
}

func (c Card) GetCost(modifiers ...CostModifier) int {
	for _, mod := range modifiers {
		c.Cost = mod(c)
	}
	return c.Cost.Cost()
}

type CardType int

const (
	CardTypePrimaryIndustry    CardType = 0
	CardTypeSecondaryIndustry  CardType = 1
	CardTypeRestaurant         CardType = 2
	CardTypeMajorEstablishment CardType = 3
	CardTypeLandMark           CardType = 4

	CardTypeBlue   = CardTypePrimaryIndustry
	CardTypeGreen  = CardTypeSecondaryIndustry
	CardTypeRed    = CardTypeRestaurant
	CardTypePurple = CardTypeMajorEstablishment
	CardTypeYellow = CardTypeLandMark
)

type CardIcon int

const (
	CardIconWheat   CardIcon = 0
	CardIconCow     CardIcon = 1
	CardIconGear    CardIcon = 2
	CardIconBread   CardIcon = 3
	CardIconFactory CardIcon = 4
	CardIconFruit   CardIcon = 5
	CardIconCup     CardIcon = 6
	CardIconTower   CardIcon = 7
)
