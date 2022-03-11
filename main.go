package main

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

var baseURL string = "https://kr.indeed.com/jobs?q=python&limit=50"

func main(){
	getPage()
}

func getPage() int {
	res, err := http.Get(baseURL)

	checkErr(err)
	checkStatusCode(res.StatusCode)

	defer res.Body.Close()

	doc,err := goquery.NewDocumentFromReader(res.Body)
	
	checkErr(err)

	// fmt.Println(doc)
	doc.Find(".pagination")

	return 0	
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