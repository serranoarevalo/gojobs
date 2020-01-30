package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type job struct {
	id       string
	title    string
	location string
	salary   string
	summary  string
}

const baseURL string = "https://kr.indeed.com/jobs?q=python&limit=50"

var jobs []job

func main() {
	totalPages := getPages()
	fmt.Println("Total pages", totalPages)
	for i := 0; i < totalPages; i++ {
		getPage(i)
	}
	fmt.Println(len(jobs))
}

func getPages() int {
	pages := 0
	res, err := http.Get(baseURL)

	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		log.Fatal("Status Code:", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		total := s.Find("a").Length()
		pages = total
	})

	return pages
}

func getPage(number int) {
	pageURL := baseURL + "&start=" + strconv.Itoa(number*50)
	res, err := http.Get(pageURL)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		log.Fatal("Status Code:", res.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".jobsearch-SerpJobCard").Each(func(index int, s *goquery.Selection) {
		jobs = append(jobs, extractJob(s))
	})
}

func extractJob(s *goquery.Selection) job {
	id, _ := s.Attr("data-jk")
	title, _ := s.Find(".title>a").Attr("title")
	title = cleanString(title)
	location := s.Find(".sjcl").Text()
	location = cleanString(location)
	salary := cleanString(s.Find(".salaryText").Text())
	summary := cleanString(s.Find(".summary").Text())
	return job{id: id, title: title, location: location, salary: salary, summary: summary}
}

func cleanString(toClean string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(toClean)), " ")
}
