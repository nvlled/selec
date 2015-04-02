package selec

import (
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"strings"
)

func find(node *html.Node, pred Pred) *html.Node {
	if node == nil {
		return nil
	}
	if pred(node) {
		return node
	}
	node = node.FirstChild
	for node != nil {
		target := find(node, pred)
		if target != nil {
			return target
		}
		node = node.NextSibling
	}
	return nil
}

func findMany(node *html.Node, pred Pred) []*html.Node {
	if node == nil {
		return nil
	}
	if pred(node) {
		return []*html.Node{node}
	} else {
		var nodes []*html.Node
		node = node.FirstChild
		for node != nil {
			subNodes := findMany(node, pred)
			nodes = concat(nodes, subNodes)
			node = node.NextSibling
		}
		return nodes
	}
}

func findAll(node *html.Node, pred Pred) []*html.Node {
	if node == nil {
		return nil
	}
	var nodes []*html.Node
	if pred(node) {
		nodes = append(nodes, node)
	}
	node = node.FirstChild
	for node != nil {
		subNodes := findAll(node, pred)
		nodes = concat(nodes, subNodes)
		node = node.NextSibling
	}
	return nodes
}

// Searches the whole node tree, including the given node
// Returns one matching node
func SelectOne(node *html.Node, preds ...Pred) *html.Node {
	for _, pred := range preds {
		node = find(node, pred)
		if node == nil {
			return nil
		}
	}
	return node
}

// Searches the whole node tree but not including the
// the subtree of the matching node.
// Returns all matching nodes.
func SelectMany(node *html.Node, preds ...Pred) []*html.Node {
	nodes := []*html.Node{node}

	for _, pred := range preds {
		var nextNodes []*html.Node
		for _, node := range nodes {
			foundNodes := findMany(node, pred)
			nextNodes = concat(nextNodes, foundNodes)
		}
		nodes = nextNodes
	}

	return nodes
}

// Searches only the children of the given node.
// Returns all matching child nodes.
func SelectChild(node *html.Node, preds ...Pred) []*html.Node {
	children := getChildren(node)
	for _, p := range preds {
		var children_ []*html.Node
		for _, child := range children {
			if p(child) {
				children_ = append(children_, child)
			}
		}
		children = children_
	}
	return children
}

// Searches the whole node tree, including the given node
// Returns all matching nodes.
func SelectAll(node *html.Node, preds ...Pred) []*html.Node {
	nodes := []*html.Node{node}

	for _, pred := range preds {
		var nextNodes []*html.Node
		for _, node := range nodes {
			foundNodes := findAll(node, pred)
			nextNodes = concat(nextNodes, foundNodes)
		}
		nodes = nextNodes
	}

	return nodes
}

func MapOne(nodes []*html.Node, preds ...Pred) []*html.Node {
	var result []*html.Node
	for _, node := range nodes {
		node_ := SelectOne(node, preds...)
		if node_ != nil {
			result = append(result, node_)
		}
	}
	return result
}

func MapAll(nodes []*html.Node, preds ...Pred) []*html.Node {
	var result []*html.Node
	for _, node := range nodes {
		nodes_ := SelectAll(node, preds...)
		if nodes != nil {
			result = append(result, nodes_...)
		}
	}
	return result
}

func Filter(nodes []*html.Node, p func(i int, node *html.Node) bool) []*html.Node {
	var nodes_ []*html.Node
	for i, node := range nodes {
		if p(i, node) {
			nodes_ = append(nodes_, node)
		}
	}
	return nodes_
}

// div, p, br adds new line
// all other tags adds space
// trimspace of each text node
func TextContent(node *html.Node) string {
	if node == nil {
		return ""
	}
	if node.Type == html.TextNode {
		var lines []string
		for _, line := range strings.Split(node.Data, "\n") {
			line = strings.TrimSpace(line)
			if line != "" {
				lines = append(lines, line)
			}
		}
		return strings.Join(lines, " ")
	}
	if node.DataAtom == atom.Code {
		s := "\n"
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			s += strings.TrimSpace(c.Data)
		}
		return s + "\n"
	}

	newline := func(c *html.Node) bool {
		return ofAtom(c, atom.P, atom.Br, atom.Div, atom.Code)
	}

	s := ""
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if newline(c) {
			s += "\n"
		}
		//s += strings.TrimSpace(TextContent(c))
		s += TextContent(c)
		if c.NextSibling != nil {
			if newline(c) {
				s += "\n"
			} else {
				s += " "
			}
		}
	}
	return s
}

func AttrVal(node *html.Node, key string) string {
	if node == nil {
		return ""
	}
	for _, attr := range node.Attr {
		if key == attr.Key {
			return attr.Val
		}
	}
	return ""
}
