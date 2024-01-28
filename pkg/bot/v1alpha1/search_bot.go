package v1alpha1

import (
	"context"
	"fmt"

	"github.com/mat285/boardgames/pkg/game/v1alpha1"
)

type HeuristicSearch struct {
	Heuristic Heuristic
	Filter    Filter
	Limiter   Limiter
}

func NewHeuristicSearch(h Heuristic, f Filter, l Limiter) (Strategy, error) {
	if h == nil {
		return nil, fmt.Errorf("Heuristic function required")
	}
	strategy := &HeuristicSearch{
		Heuristic: h,
		Filter:    FilterOrDefault(f),
		Limiter:   LimiterOrDefault(l),
	}
	return strategy, nil
}

func (hs *HeuristicSearch) ChooseMove(ctx context.Context, state v1alpha1.StateData) (v1alpha1.Move, error) {
	nodes, err := hs.Search(state)
	if err != nil {
		return nil, err
	}
	return nodes[0].Move, nil
}

func (hs *HeuristicSearch) Search(state v1alpha1.StateData) ([]*Node, error) {
	tree := NewTree(state, hs.Heuristic, hs.Filter)

	for !hs.Limiter(tree.Size, tree.Depth) {
		err := tree.ExpandLeaves()
		if err != nil {
			return nil, err
		}
	}

	pid, err := state.CurrentPlayer()
	if err != nil {
		return nil, err
	}

	err = tree.BackpropogateScores(pid)
	if err != nil {
		return nil, err
	}

	explored := make([]*Node, 0, len(tree.Root.Children))

	for _, child := range tree.Root.Children {
		if child == nil {
			continue
		}
		explored = append(explored, child)
	}
	SortNodes(explored)
	return explored, nil
}
