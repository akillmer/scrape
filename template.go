package scrape

import (
	"bytes"
	"io"
	"strings"

	strip "github.com/grokify/html-strip-tags-go"
	"golang.org/x/net/html"
)

var textToken = html.TextToken.String()

// Template for scraping
type Template struct {
	targets    []*templateNode
	validators map[string]func(v string) bool
}

func renderNodeText(n *html.Node) (string, error) {
	buf := bytes.NewBufferString("")
	if err := html.Render(buf, n); err != nil {
		return "", err
	}
	return strip.StripTags(buf.String()), nil
}

// NewTemplate with provided template string
func NewTemplate(template string) *Template {
	tmpl := &Template{
		validators: make(map[string]func(v string) bool),
	}

	var cursor *templateNode

	it := html.NewTokenizerFragment(strings.NewReader(template), "")
	for t := it.Next(); t != html.ErrorToken; t = it.Next() {
		token := it.Token()

		if t == html.StartTagToken || t == html.SelfClosingTagToken {
			node := newTemplateNode(&token)
			node.parent = cursor
			cursor = node
		}

		if t.String() == textToken {
			text := strings.TrimSpace(token.Data)
			if len(text) > 0 {
				cursor.ParseAttribute(textToken, text)
			}
		}

		if t == html.EndTagToken || t == html.SelfClosingTagToken {
			if len(cursor.scrapeAttr) > 0 {
				tmpl.targets = append(tmpl.targets, cursor)
			}
			cursor = cursor.parent
		}
	}

	return tmpl
}

// Validator for defining functions used within a template
func (tmpl *Template) Validator(name string, fn func(v string) bool) *Template {
	tmpl.validators[name] = fn
	return tmpl
}

// Scrape the provided Reader interface
func (tmpl *Template) Scrape(r io.Reader) (map[string][]string, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	results := make(map[string][]string)

	var traverse func(n *html.Node)
	traverse = func(n *html.Node) {
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			for _, target := range tmpl.targets {
				if r := tmpl.scrapeNode(target, child); r != nil {
					for k, v := range r {
						results[k] = append(results[k], v)
					}
				}
			}
			traverse(child)
		}
	}

	traverse(doc)

	return results, nil
}

func (tmpl *Template) scrapeNode(tn *templateNode, dn *html.Node) map[string]string {
	m := make(map[string]string)

	for target := tn; target != nil; target = target.parent {
		if target.tag != dn.Data {
			return nil
		}

		nodeAttrs := func() map[string]string {
			m := make(map[string]string)
			for _, attr := range dn.Attr {
				m[attr.Key] = attr.Val
			}
			return m
		}()

		for k, v := range target.requireAttr {
			if nodeValue, exists := nodeAttrs[k]; !exists || v != nodeValue {
				return nil
			}
		}

		for k, v := range target.validateAttr {
			if fn, exists := tmpl.validators[v]; exists == false {
				panic("validator `" + v + "()` does not exist within the template")
			} else {
				if value, exists := nodeAttrs[k]; exists == false {
					return nil
				} else if fn(value) == false {
					return nil
				}
			}
		}

		for k, v := range target.scrapeAttr {
			if k == textToken {
				m[v], _ = renderNodeText(dn)
			} else {
				m[v] = nodeAttrs[k]
			}
		}

		dn = dn.Parent
	}

	return m
}
