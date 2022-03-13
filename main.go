package main

import (
	"os"
	"strings"

	cleanstring "github.com/baejoonsoo/webScraper/cleanString"
	"github.com/baejoonsoo/webScraper/scraper"
	"github.com/labstack/echo"
)

var fileName string = "jobs.csv"

func handleHome(c echo.Context) error{
	return c.File("home.html")
}

func handleScrape(c echo.Context) error{
	defer os.Remove(fileName)
	term := strings.ToLower(cleanstring.CleanString(c.FormValue("term")))
	scraper.Scrape(term)
	return c.Attachment(fileName,fileName)
}

func main(){
	e := echo.New()

	e.GET("/",handleHome)
	e.POST("/crape",handleScrape)

	e.Logger.Fatal(e.Start(":8000"))
}