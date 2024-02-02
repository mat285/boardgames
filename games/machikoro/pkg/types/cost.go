package types

type Cost struct {
	Values [3]int
	C      int
}

func NewCost(c int) Cost {
	return Cost{
		Values: [3]int{c, c, c},
		C:      c,
	}
}

func NewLandmarkCost(first, second, third int) Cost {
	return Cost{
		Values: [3]int{first, second, third},
		C:      first,
	}
}

func (c Cost) Cost() int {
	return c.C
}

type CostModifier func(Card) Cost

func LandmarkModifier(num int) CostModifier {
	return func(c Card) Cost {
		if c.Type != CardTypeLandMark {
			return c.Cost
		}
		cost := c.Cost
		cost.C = cost.Values[num]
		return cost
	}
}

func LaunchPadModifier() CostModifier {
	return func(c Card) Cost {
		if c.Name != "Launchpad" {
			return c.Cost
		}
		cost := c.Cost.Values
		return NewLandmarkCost(cost[0]/2, cost[1]/2, cost[2]/2)
	}
}
