package v1alpha1

import (
	"sort"

	"github.com/blend/go-sdk/uuid"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
)

type ScoredState interface {
	GetState() v1alpha1.StateData
	GetScore() int
}

type Level []*Node

type Node struct {
	State    v1alpha1.StateData
	Move     v1alpha1.Move
	Score    int
	Parent   *Node
	Children []*Node

	Tree *Tree
}

func (n *Node) AddChild(children ...*Node) {
	n.Children = append(n.Children, children...)
	if len(n.Children) > 0 {
		delete(n.Tree.Leaves, n)
	}
}

func (n *Node) FromPlayerMove(id uuid.UUID) (bool, error) {
	// We need to check that in the previous state i.e.
	// the parent of this node was our current player
	// otherwise this state was caused by another player
	// so we don't have control over getting here perfectly
	if n.Parent == nil {
		return false, nil
	}
	curr, err := n.Parent.State.CurrentPlayer()
	if err != nil {
		return false, err
	}
	return id.Equal(curr), nil
}

func (n *Node) GetState() v1alpha1.StateData {
	return n.State
}

func (n *Node) GetScore() int {
	return n.Score
}

func SortNodes(nodes []*Node) {
	// sorts in descending order
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Score > nodes[j].Score
	})
}
