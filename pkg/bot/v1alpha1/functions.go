package v1alpha1

import "github.com/mat285/boardgames/pkg/game/v1alpha1"

type Heuristic func(v1alpha1.StateData) int

type Filter func([]*Node) []*Node

type Limiter func(size, depth int) (stop bool)

func FilterOrDefault(f Filter) Filter {
	if f == nil {
		return DefaultFilter
	}
	return f
}

func LimiterOrDefault(l Limiter) Limiter {
	if l == nil {
		return DefaultLimiter
	}
	return l
}

func DefaultFilter(nodes []*Node) []*Node {
	SortNodes(nodes)
	limit := 5
	if len(nodes) < limit {
		limit = len(nodes)
	}
	return nodes[:limit]
}

func DefaultLimiter(size, depth int) (stop bool) {
	stop = size > (1 << 20)
	return
}
