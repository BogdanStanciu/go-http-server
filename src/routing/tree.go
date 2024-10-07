package routing

import (
	"fmt"
	"go-http-server/src/http"
	"log"
)

type Node struct {
	Name    string
	Child   []*Node
	IsLeaf  bool
	Handler http.RouteHandlerFunction
}

// Init a empty tree
func InitTree() *Node {
	return &Node{Name: "", Child: make([]*Node, 0), IsLeaf: false, Handler: nil}
}

// Add a node in the tree.
func (tree *Node) AddNode(path []string, handler http.RouteHandlerFunction) {
	if len(path) > 0 {
		var i int
		for i = 0; i < len(tree.Child); i++ {
			if tree.Child[i].Name == path[0] {
				break
			}
		}
		if i == len(tree.Child) {
			tree.Child = append(tree.Child, &Node{Name: path[0], Child: make([]*Node, 0), IsLeaf: false})
		}
		tree.Child[i].AddNode(path[1:], handler)
	}
	tree.Handler = (handler)
	tree.IsLeaf = true
}

// Search for a path in the tree structure
func (tree *Node) Search(path []string) *Node {
	if len(path) == 0 && len(tree.Child) == 0 {
		return tree
	}

	var i int
	for i = 0; i < len(tree.Child); i++ {
		if tree.Child[i].Name == path[0] || tree.Child[i].Name == "*" {
			break
		}
	}

	if i == len(tree.Child) {
		return nil
	}

	return tree.Child[i].Search(path[1:])
}

// Print all route of a given tree
func (tree *Node) PrintTree(base string) {
	if len(tree.Child) == 0 {
		log.Printf("[Router] Mapped %s%s\n", base, tree.Name)
		return
	}

	for i := range tree.Child {
		tree.Child[i].PrintTree(fmt.Sprintf("%s%s/", base, tree.Name))
	}
}
