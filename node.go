package scrape

import (
	"strings"

	"golang.org/x/net/html"
)

type templateNode struct {
	parent       *templateNode
	tag          string
	requireAttr  map[string]string
	scrapeAttr   map[string]string
	validateAttr map[string]string
}

func newTemplateNode(token *html.Token) *templateNode {
	node := &templateNode{
		tag:          token.Data,
		requireAttr:  make(map[string]string),
		scrapeAttr:   make(map[string]string),
		validateAttr: make(map[string]string),
	}

	for _, attr := range token.Attr {
		node.ParseAttribute(attr.Key, attr.Val)
	}

	return node
}

func (tn *templateNode) ParseAttribute(key, val string) {
	if strings.HasPrefix(val, "{{") && strings.HasSuffix(val, "}}") {
		for _, v := range strings.Split(val[2:len(val)-2], "|") {
			if strings.HasSuffix(v, "()") {
				tn.validateAttr[key] = v[:len(v)-2]
			} else {
				tn.scrapeAttr[key] = v
			}
		}
	} else {
		tn.requireAttr[key] = val
	}
}
