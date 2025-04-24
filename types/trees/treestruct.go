package trees

import "fmt"

type BiTreeNode struct {
	Val   int
	Left  *BiTreeNode
	Right *BiTreeNode
}

func (t *BiTreeNode) String() string {
	if t == nil {
		return "nil"
	}
	return fmt.Sprintf("{Val: %d, Left: %v, Right: %v}", t.Val, t.Left, t.Right)
}

type MultiChildTreeNode struct {
	Val       int                   `json:"val"`
	Children  []*MultiChildTreeNode `json:"children"`
	IsVisited bool                  `json:"isVisited"`
	Metadata  TreeMetadata          `json:"metadata"`
}

func (t *MultiChildTreeNode) String() string {
	if t == nil {
		return "nil"
	}
	return fmt.Sprintf("{Val: %d, Children: %v}", t.Val, t.Children)
}

type TreeMetadata struct {
	Label string
	Color string
	Depth int
}
