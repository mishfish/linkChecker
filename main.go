package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/html"
)

func getHTML(url string) (body []byte, e error) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
func getLinks(body []byte) (links []string, e error) {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					links = append(links, a.Val)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return links, nil
}

func check(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	if resp.StatusCode == 200 {
		fmt.Println("OK", url)
	} else {
		fmt.Println("ERR", url)
	}
}

func main() {
	var err error
	var content []byte
	url := flag.String("url", "", "url adress")
	filepath := flag.String("filepath", "", "file adress")
	flag.Parse()
	if flag.NArg() != 0 {
		panic("too many args")
	}
	//@todo убрать при первой возможности
	if *url != "" {
		content, err = getHTML(*url)
	} else {
		content, err = ioutil.ReadFile(*filepath)
	}
	if err != nil {
		panic(err)
	}
	links, err := getLinks(content)
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	for _, link := range links {
		wg.Add(1)
		go check(link, &wg)
	}
	wg.Wait()
}
