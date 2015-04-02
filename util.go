package selec

import (
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
)

func concat(nodes1, nodes2 []*html.Node) []*html.Node {
	if nodes1 == nil {
		return nodes2
	} else if nodes2 == nil {
		return nodes1
	}
	return append(nodes1, nodes2...)
}

func getChildren(node *html.Node) []*html.Node {
	var children []*html.Node
	child := node.FirstChild
	for child != nil {
		if child.Type != html.TextNode {
			children = append(children, child)
		}
		child = child.NextSibling
	}
	return children
}

func getChild(index int, node *html.Node) *html.Node {
	node = node.FirstChild
	for i := 0; node != nil; {
		if i == index {
			return node
		}
		if node.Type != html.TextNode {
			i++
		}
		node = node.NextSibling
	}
	return nil
}

// returns the index of the node with respect to its parent
func indexOf(node *html.Node) int {
	parent := node.Parent
	if parent == nil {
		return -1
	}
	child := parent.FirstChild
	i := 0
	for {
		if node == child {
			return i
		}
		if child.Type != html.TextNode {
			i++
		}
		child = child.NextSibling
	}
	return -1
}

func ofAtom(node *html.Node, atoms ...atom.Atom) bool {
	for _, a := range atoms {
		if node.DataAtom == a {
			return true
		}
	}
	return false
}
