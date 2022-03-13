package scraper

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	cleanstring "github.com/baejoonsoo/webScraper/cleanString"
)

// 크롤링 짧은 시간내에 여러번 계속하면 차단 먹음...
// 조심하기...

type extractedJob struct {
	id string
	title string
	location string
	salary string
	summary string
}

// Scrape indeed by a term
func Scrape(term string){
	baseURL := "https://kr.indeed.com/jobs?q="+term+"&limit=50"
	LastpageURL := "https://kr.indeed.com/jobs?q="+term+"&limit=50&start=9999"

	channel := make(chan []extractedJob)
	jobs := []extractedJob{}
	totalPages :=	getPages(LastpageURL)

	
	for i := 0;i<totalPages; i++{
		go getPage(i, baseURL, channel)
	}
	
	for i := 0;i<totalPages; i++{
	extractedJobs := <-channel
	jobs = append(jobs, extractedJobs...)
	}

	writeJobs(jobs)
	fmt.Println("Done, extracted", len(jobs))

}

func getPage(page int, baseURL string, mainChannel chan<- []extractedJob) {
	channel := make(chan extractedJob)

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
		go extractJob(card, channel)
	})

	for i:=0; i<searchCards.Length(); i++{
		job := <-channel
		jobs = append(jobs, job)
	}

	mainChannel <- jobs
}

func extractJob(card *goquery.Selection, channel chan<- extractedJob) {
	id,_ := card.Attr("data-jk")		
	title:= cleanstring.CleanString(card.Find(".jobTitle>span").Text())
	location := cleanstring.CleanString(card.Find(".companyLocation").Text())
	salary := cleanstring.CleanString(card.Find(".salary-snippet").Text())
	summary := cleanstring.CleanString(card.Find(".job-snippet").Text())
	
	channel <- extractedJob{
		id: id,
		title: title,
		location: location,
		salary: salary,
		summary: summary,
	}
}




func getPages(LastpageURL string) int {
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

func writeJobs(jobs []extractedJob){
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"Link", "Title", "Location", "Salary", "Summary"}

	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs{
		jobSlice := []string{"https://kr.indeed.com/viewjob?jk="+job.id, job.title, job.location, job.salary, job.summary}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
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