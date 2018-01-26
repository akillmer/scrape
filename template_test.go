package scrape

import (
	"os"
	"testing"

	"golang.org/x/net/html"
)

func TestRenderNodeText(t *testing.T) {
	f, err := os.Open("./test.html")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	doc, err := html.Parse(f)
	if err != nil {
		t.Fatal(err)
	}

	var traverse func(n *html.Node)
	var span *html.Node

	traverse = func(n *html.Node) {
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			if child.Data == "span" {
				span = child
				break
			}

			if span == nil {
				traverse(child)
			}
		}
	}

	traverse(doc)

	if span == nil {
		t.Fatal("could not find <span> node within \"test.html\" document")
	}

	spanText, err := renderNodeText(span)
	if err != nil {
		t.Fatal(err)
	}

	if spanText != "Super grumpy cat!" {
		t.Fatalf("expected \"Super grumpy cat!\", got \"%s\"", spanText)
	}
}

func TestTemplateScrape(t *testing.T) {
	f, err := os.Open("./test.html")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	tmpl := NewTemplate(
		`<div class="gallery">
			<img src="{{imgSrc}}"/>
			<span>{{imgText}}</span>
		</div>
	`)

	_, err = tmpl.Scrape(f)
	if err != nil {
		t.Fatal(err)
	}
}
