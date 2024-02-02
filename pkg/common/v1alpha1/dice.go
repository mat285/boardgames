package v1alpha1

import "math/rand"

type Die interface {
	Faces() int
	Value(int) int
}

func Roll(d Die) int {
	return d.Value(rand.Intn(d.Faces()))
}

func Sum(ds ...Die) int {
	s := 0
	for _, d := range ds {
		s += Roll(d)
	}
	return s
}

func DieN(n int) Die {
	return ContinuousDie{Num: n}
}

type ContinuousDie struct {
	Num int
}

func (cd ContinuousDie) Faces() int {
	return cd.Num
}

func (cd ContinuousDie) Value(f int) int {
	if f >= 0 && f < cd.Num {
		return f + 1
	}
	return 0
}
