package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id string
	location string
	title string
	salary string
	summary string
}

var baseURL string = "https://kr.indeed.com/jobs?q=python&limit=50"

var LastpageURL string = "https://kr.indeed.com/jobs?q=python&limit=50&start=9999"

func main(){
	totalPages :=	getPages()
	
	for i := 0;i<totalPages; i++{
		getPage(i)

	}
}

func getPage(page int){
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
		id,_ := card.Attr("data-jk")		
		title:= card.Find(".jobTitle>span").Text()
		location := card.Find(".companyLocation").Text()
		summary := card.Find(".job-snippet").Text()
		
		fmt.Println(id, title, location, summary)
	})
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

/*
func getPages() int {
pages := 0
res, err := http.Get(baseURL)

checkErr(err)
checkStatusCode(res.StatusCode)

defer res.Body.Close()

doc,err := goquery.NewDocumentFromReader(res.Body)

checkErr(err)

// fmt.Println(doc)
doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
	pages = s.Find("a").Length()
})

return pages	
}
*/
