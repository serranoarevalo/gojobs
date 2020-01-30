package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

const baseURL string = "https://kr.indeed.com/jobs?q=python&limit=50"

func main() {
	totalPages := getPages()
	fmt.Println(totalPages)
}

func getPages() int {
	c := colly.NewCollector()
	pages := 0
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL.String())
	})
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited:", r.Request.URL.String())
	})

	c.OnHTML(".pagination", func(e *colly.HTMLElement) {
		e.DOM.Each(func(i int, s *goquery.Selection) {
			total := s.Find("a").Length()
			pages = total
		})
	})
	c.Visit(baseURL)
	return pages
}
