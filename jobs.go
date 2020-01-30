package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
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

func main() {
	var jobs []job
	c := make(chan []job)

	totalPages := getPages()
	fmt.Println("Extracted", totalPages, "pages")

	for i := 0; i < totalPages; i++ {
		go getPage(i, c)
	}

	for i := 0; i < totalPages; i++ {
		pageJobs := <-c
		jobs = append(jobs, pageJobs...)
	}

	fmt.Println("Writting", len(jobs), "jobs")
	writeJobs(jobs)
	fmt.Println("Done")
}

func getPage(number int, mainChannel chan []job) {
	var jobs []job
	pageURL := baseURL + "&start=" + strconv.Itoa(number*50)
	fmt.Println("Requesting", pageURL)
	res, err := http.Get(pageURL)
	checkError(err)
	checkStatusCode(res)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkError(err)

	innerChannel := make(chan job)

	searchCards := doc.Find(".jobsearch-SerpJobCard")

	searchCards.Each(func(index int, s *goquery.Selection) {
		go extractJob(s, innerChannel)
	})

	for i := 0; i < searchCards.Length(); i++ {
		extracted := <-innerChannel
		jobs = append(jobs, extracted)
	}
	mainChannel <- jobs
}

func extractJob(s *goquery.Selection, c chan job) {
	id, _ := s.Attr("data-jk")
	title, _ := s.Find(".title>a").Attr("title")
	title = cleanString(title)
	location := s.Find(".sjcl").Text()
	location = cleanString(location)
	salary := cleanString(s.Find(".salaryText").Text())
	summary := cleanString(s.Find(".summary").Text())
	c <- job{id: "https://www.indeed.com/viewjob?jk=" + id, title: title, location: location, salary: salary, summary: summary}
}

func getPages() int {
	pages := 0
	res, err := http.Get(baseURL)
	checkError(err)
	checkStatusCode(res)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkError(err)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		total := s.Find("a").Length()
		pages = total
	})

	return pages
}

func cleanString(toClean string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(toClean)), " ")
}

func writeJobs(jobs []job) {
	file, err := os.Create("jobs.csv")
	checkError(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"apply", "title", "location", "salary", "summary"}
	writeErr := w.Write(headers)
	checkError(writeErr)

	for _, job := range jobs {
		jobCSV := []string{job.id, job.title, job.location, job.salary, job.summary}
		writeErr := w.Write(jobCSV)
		checkError(writeErr)
	}

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkStatusCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatal("Status Code:", res.StatusCode)
	}
}
