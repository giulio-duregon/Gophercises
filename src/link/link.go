package link

import (
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"log"
	"os"
)

type Link struct {
	href, text string
}

func parseCli() (fpath *string) {
	fpath = flag.String("file", "static/link/ex2.html", "The file path for the HTML document you want to parse")
	return
}

func dfs(n *html.Node) ([]*Link, bool) {
	if n == nil {
		return nil, false
	}

	var links []*Link
	if n.DataAtom == atom.A {
		link := Link{}

		for _, at := range n.Attr {
			if at.Key == "href" {
				link.href = at.Val
			}
		}

		link.text = n.FirstChild.Data
		links = append(links, &link)
	}
	if n.FirstChild != nil {
		ls, ok := dfs(n.FirstChild)
		if ok {
			links = append(links, ls...)
		}

	}
	if n.NextSibling != nil {
		ls, ok := dfs(n.NextSibling)
		if ok {
			links = append(links, ls...)
		}
	}
	return links, true
}

func LinkProgram() {
	filePath := parseCli()

	// Open file, defer closing, handle errors
	f, err := os.Open(*filePath)
	if err != nil {
		log.Fatalf("Couldn't read file %s, %v\n", *filePath, err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("Couldnt close file %s, %v\n", *filePath, err)
		}
	}(f)

	// Parse HTML, Handle Errors
	head, err := html.Parse(f)
	if err != nil {
		log.Fatalf("Couldn't parse HTML, %v", err)
	}
	// Run DFS search ok links, ok == false means none found
	out, ok := dfs(head)
	if ok {
		for i, l := range out {
			fmt.Printf("%v: %+v\n", i, *l)
		}
	}
}
