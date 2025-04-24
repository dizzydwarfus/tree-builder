package main

import (
	"time"

	"github.com/dizzydwarfus/tree-builder/graph/tool"
	"github.com/dizzydwarfus/tree-builder/internal/shared"
	"github.com/dizzydwarfus/tree-builder/internal/treetraversal"
	"github.com/dizzydwarfus/tree-builder/types/trees"
)

func main() {
	// always start with 1 as root node
	start := time.Now()
	tree := &trees.MultiChildTreeNode{
		Val:       1,
		Children:  []*trees.MultiChildTreeNode{},
		IsVisited: false,
		Metadata: trees.TreeMetadata{
			Label: "root",
			Color: shared.Colors[0],
			Depth: 0,
		},
	}
	value := 2 // need to refactor to remove dependency on value variable in TreeBuilder
	var counter *int = &value
	treeInput := []int{2, 2, 2, 2} // number of children per node from second level onwards
	treetraversal.TreeBuilder(tree, treeInput, counter, 1)
	treetraversal.ShowJSONTree(tree)
	shared.Yellow("TreeBuilder BFS took: %v\n", time.Since(start))

	start = time.Now()
	treeGraph := tool.CreateDotFile("testtree", "testtree")
	tool.CreateTreeGraph(tree, treeGraph)
	tool.CloseDotFile(treeGraph)
	tool.CreateGraph(treeGraph, "gif", false)
	shared.Yellow("Graph creation took: %v\n", time.Since(start))

	start = time.Now()
	bfsDepth := treetraversal.BfsMultiChild(tree)
	shared.Yellow("BFS took: %v\n", time.Since(start))

	start = time.Now()
	dfsDepth := 0
	dfsMaxDepth := 0
	treetraversal.DfsMultiChild(tree, &dfsDepth, &dfsMaxDepth)
	shared.Yellow("DFS took: %v\n", time.Since(start))

	shared.Green("BFS Depth: %v\n", bfsDepth)
	shared.Green("DFS Depth: %v\n", dfsMaxDepth)
}
