package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/html"
)

func getHTML(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body
}
func getLinks(body []byte) []string {
	doc, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
	var links []string
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
	return links
}

func check(url string, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	if resp.StatusCode == 200 {
		fmt.Println("OK", url)
	} else {
		fmt.Println("error", url)
	}
}

func main() {
	html := getHTML("http://google.com")
	links := getLinks(html)
	var wg sync.WaitGroup
	for i, link := range links {
		wg.Add(1)
		println("start", i, link)
		go check(link, &wg)
	}
	wg.Wait()
}
