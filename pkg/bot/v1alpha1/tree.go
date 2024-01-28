package v1alpha1

import (
	"fmt"

	"github.com/blend/go-sdk/uuid"
	"github.com/mat285/boardgames/pkg/game/v1alpha1"
)

type Tree struct {
	Root   *Node
	Leaves map[*Node]bool

	Size  int
	Depth int

	Heuristic Heuristic
	Filter    Filter
}

func NewTree(state v1alpha1.StateData, h Heuristic, f Filter) *Tree {
	root := &Node{
		State: state,
		Score: h(state),
	}
	tree := &Tree{
		Root: root,
		Leaves: map[*Node]bool{
			root: true,
		},

		Size:  1,
		Depth: 1,

		Heuristic: h,
		Filter:    f,
	}
	root.Tree = tree
	return tree
}

func (t *Tree) ExpandLeaves() error {
	for leaf := range t.Leaves {
		err := t.ExpandNode(leaf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Tree) ExpandNode(node *Node) error {
	possible, err := node.State.ValidMoves()
	if err != nil {
		return err
	}

	children := make([]*Node, 0, len(possible))
	for _, move := range possible {
		res, err := move.Apply(node.State)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if !res.Valid {
			continue
		}
		child := &Node{
			State:  res.State,
			Move:   move,
			Score:  t.Heuristic(res.State),
			Parent: node,
			Tree:   node.Tree,
		}
		children = append(children, child)
		node.Tree.Leaves[child] = true
		t.Size++
	}
	node.AddChild(t.Filter(children)...)
	t.Depth++
	return nil
}

func (t *Tree) BackpropogateScores(id uuid.UUID) error {
	candidates := make(map[*Node]bool)
	next := make(map[*Node]bool)
	processed := make(map[*Node]bool)

	for leaf := range t.Leaves {
		if leaf == nil || leaf.Parent == nil {
			continue
		}
		candidates[leaf.Parent] = true
		processed[leaf] = true
	}

	// leaves nodes are soley determined by their hueristic score
	// nodes further up are either the max of their children
	// if we are the current player, otherwise it's the weighted avergae
	// to account for other players choices

	for len(candidates) > 0 {
		for node := range candidates {
			if processed[node] {
				continue
			}
			processed[node] = true
			us, err := node.FromPlayerMove(id)
			if err != nil {
				return err
			}
			if us {
				node.Score = maxNode(node.Children).Score
			} else {
				node.Score = averageNodeScore(node.Children)
			}
			if node.Parent != nil {
				next[node.Parent] = true
			}
		}
		candidates = next
		next = make(map[*Node]bool)
	}

	return nil
}

func averageNodeScore(nodes []*Node) int {
	if len(nodes) == 0 {
		return 0
	}
	sum := 0
	for _, node := range nodes {
		sum += node.Score
	}
	return sum / len(nodes)
}

func maxNode(nodes []*Node) *Node {
	if len(nodes) == 0 {
		return nil
	}
	max := nodes[0]
	for i := 1; i < len(nodes); i++ {
		if nodes[i].Score > max.Score {
			max = nodes[i]
		}
	}
	return max
}
