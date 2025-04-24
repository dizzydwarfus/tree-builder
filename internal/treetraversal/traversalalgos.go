package treetraversal

import (
	"github.com/dizzydwarfus/tree-builder/types/trees"
)

func BfsSimple(root *trees.BiTreeNode) int {
	if root == nil {
		return 0
	}

	nodeList := []*trees.BiTreeNode{root}
	depth := 0

	for len(nodeList) > 0 {

		for _, node := range nodeList {
			// fmt.Println("Before Slicing", nodeList)
			nodeList = nodeList[1:]
			// fmt.Println("After Slicing", nodeList)

			// fmt.Println("Current node value: ", node.Val)

			if node.Left != nil {
				// fmt.Println("Appending left node: ", node.Left)
				nodeList = append(nodeList, node.Left)
			}
			if node.Right != nil {
				// fmt.Println("Appending right node: ", node.Right)
				nodeList = append(nodeList, node.Right)
			}
			// fmt.Println()
		}

		depth++
	}

	return depth
}

func BfsMultiChild(root *trees.MultiChildTreeNode) int {
	if root == nil {
		return 0
	}

	nodeList := []*trees.MultiChildTreeNode{root}
	depth := 0

	for len(nodeList) > 0 {
		for _, node := range nodeList {
			// fmt.Println("Before Slicing", nodeList)
			nodeList = nodeList[1:]
			// fmt.Println("After Slicing", nodeList)

			// fmt.Println("Current node value: ", node.Val)

			if node.Children != nil {
				// fmt.Println("Appending children nodes: ", node.Children)
				nodeList = append(nodeList, node.Children...)
			}
			// fmt.Println()
		}
		depth++
	}
	return depth
}

func DfsMultiChild(node *trees.MultiChildTreeNode, depth *int, maxDepth *int) {
	// condition should be when node == nil
	// recursively call current func for each child
	// maybe need isVisited bool?
	//TODO: if tree is given without IsVisited field, how to extend tree into a new struct or use interface
	// problem is function recursion takes in *MultiChildTreeNode and if new struct, cannot default call same function recursively
	if node == nil {
		return
	}
	// shared.Faint("Setting IsVisited to True\n")
	node.IsVisited = true
	*depth++
	if *depth > *maxDepth {
		// shared.Faint("Setting maxDepth %v to new max depth: %v\n", *maxDepth, *depth)
		*maxDepth = *depth
	}
	if len(node.Children) > 0 {
		// shared.Faint("Looping through node.Children at node.Val: %v\n", node.Val)
		for _, child := range node.Children {
			// shared.Faint("Recursing to child.Val: %v\n", child.Val)
			DfsMultiChild(child, depth, maxDepth)
		}
	}
	// shared.Faint("Backtracking...\n")
	*depth--
}
