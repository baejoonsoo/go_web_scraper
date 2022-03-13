package main

import (
	"fmt"
	"strings"

	"github.com/baejoonsoo/webScraper/scraper"
	"github.com/labstack/echo"
)

func handleHome(c echo.Context) error{
	return c.File("home.html")
}

func handleScrape(c echo.Context) error{
	term := strings.ToLower(scraper.CleanString(c.FormValue("term")))
	fmt.Println(term)
	return nil
}

func main(){
	e := echo.New()

	e.GET("/",handleHome)
	e.POST("/crape",handleScrape)

	e.Logger.Fatal(e.Start(":8000"))
}