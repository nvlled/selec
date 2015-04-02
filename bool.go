package selec

import (
	"code.google.com/p/go.net/html"
	"strings"
)

type Pred func(node *html.Node) bool

func And(preds ...Pred) Pred {
	return func(node *html.Node) bool {
		for _, holds := range preds {
			if !holds(node) {
				return false
			}
		}
		return true
	}
}

func Or(preds ...Pred) Pred {
	return func(node *html.Node) bool {
		for _, holds := range preds {
			if holds(node) {
				return true
			}
		}
		return false
	}
}

func Not(pred Pred) Pred {
	return func(node *html.Node) bool {
		return !pred(node)
	}
}

func WithAttr(key string, fn func(val string) bool) Pred {
	return func(node *html.Node) bool {
		for _, attr := range node.Attr {
			if key == attr.Key && fn(attr.Val) {
				return true
			}
		}
		return false
	}
}

func HasPrefix(pref string) func(string) bool {
	return func(val string) bool {
		return strings.HasPrefix(val, pref)
	}
}

func HasSubstr(sub string) func(string) bool {
	return func(val string) bool {
		return strings.Contains(val, sub)
	}
}

func Attr(key, val string) Pred {
	return WithAttr(key, func(val_ string) bool {
		return val == val_
	})
}

func Id(name string) Pred {
	return WithAttr("id", func(classNames string) bool {
		for _, name_ := range strings.Fields(classNames) {
			if name == name_ {
				return true
			}
		}
		return false
	})
}

func Class(name string) Pred {
	return WithAttr("class", func(classNames string) bool {
		for _, name_ := range strings.Fields(classNames) {
			if name == name_ {
				return true
			}
		}
		return false
	})
}

func AttrOnly(key, val string) Pred {
	return func(node *html.Node) bool {
		if len(node.Attr) != 1 {
			return false
		}
		attr := node.Attr[0]
		return key == attr.Key &&
			val == attr.Val
	}
}

func Tag(name string) Pred {
	return func(node *html.Node) bool {
		if node.Type == html.ElementNode {
			return name == node.Data
		}
		return false
	}
}

func TagAttr(name, key, val string) Pred {
	return And(Tag(name), Attr(key, val))
}

func TagAttrOnly(name, key, val string) Pred {
	return And(Tag(name), AttrOnly(key, val))
}

func Text(node *html.Node) bool {
	return node.Type == html.TextNode
}

func Nth(n int) Pred {
	return func(node *html.Node) bool {
		parent := node.Parent
		if parent == nil {
			return false
		}
		return getChild(n, parent) == node
	}
}

func EveryNth(n int) Pred {
	return func(node *html.Node) bool {
		i := indexOf(node)
		return i >= 0 && (i+1)%n == 0
	}
}

func Last(p Pred) Pred {
	return func(node *html.Node) bool {
		parent := node.Parent
		if parent == nil {
			return false
		}
		for c := parent.LastChild; c != nil; c = c.PrevSibling {
			if p(c) {
				return c == node
			}
		}
		return false
	}
}

// node.LastChild actually includes text nodes
//func LastChild(node *html.Node) bool {
//	parent := node.Parent
//	return parent != nil && parent.LastChild == node
//}
