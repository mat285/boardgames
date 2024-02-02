package v1alpha1

type IntRange struct {
	Low  int
	High int
}

func NewIntRange(low, high int) IntRange {
	return IntRange{
		Low:  low,
		High: high,
	}
}

func (ir IntRange) Includes(i int) bool {
	return i >= ir.Low && i <= ir.High
}
