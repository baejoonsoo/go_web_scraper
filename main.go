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

// 크롤링 여러번 계속하면 차단 먹음...
// 조심하기...

type extractedJob struct {
	id string
	title string
	location string
	salary string
	summary string
}

var baseURL string = "https://kr.indeed.com/jobs?q=python&limit=50"

var LastpageURL string = "https://kr.indeed.com/jobs?q=python&limit=50&start=9999"

func main(){
	jobs := []extractedJob{}
	totalPages :=	getPages()

	
	for i := 0;i<totalPages; i++{
		extractedJobs := getPage(i)
		jobs = append(jobs, extractedJobs...)
	}
	// fmt.Println(jobs)
	writeJobs(jobs)
}

func writeJobs(jobs []extractedJob){
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"ID", "Title", "Location", "Salary", "Summary"}

	wErr := w.Write(headers)
	checkErr(wErr)

}

func getPage(page int) []extractedJob{
	jobs := []extractedJob{}
	pageURL := baseURL + "&start="+ strconv.Itoa(page*50)
	fmt.Println("Requesting",pageURL)

	res, err := http.Get(pageURL)

	checkErr(err)
	checkStatusCode(res.StatusCode)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".tapItem")

	searchCards.Each(func(i int, card *goquery.Selection) {
		job := extractJob(card)
		jobs = append(jobs, job)
	})

	return jobs
}

func extractJob(card *goquery.Selection) extractedJob{
	id,_ := card.Attr("data-jk")		
	title:= cleanString(card.Find(".jobTitle>span").Text())
	location := cleanString(card.Find(".companyLocation").Text())
	salary := cleanString(card.Find(".salary-snippet").Text())
	summary := cleanString(card.Find(".job-snippet").Text())
	
	return extractedJob{
		id: id,
		title: title,
		location: location,
		salary: salary,
		summary: summary,
	}
}

func cleanString(str string) string{
	return strings.Join(strings.Fields(strings.TrimSpace(str))," ")
}

func getPages() int {
	pages := 0
	res, err := http.Get(LastpageURL)
	checkErr(err)
	checkStatusCode(res.StatusCode)
	
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)
	
	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages, err = strconv.Atoi(s.Find("b").Text())
	})
	return pages
}

func checkErr(err error){
	if err != nil{
		log.Fatalln(err)
	}
}

func checkStatusCode(StatusCode int){
	if StatusCode != 200{
		log.Fatalln("Request failed with Status : ", StatusCode)
	}
}