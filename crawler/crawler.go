package crawler

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

func starColly() {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.113 Safari/537.36"),
	)

	c.OnHTML("h3", func(e *colly.HTMLElement) {
		fmt.Printf("section %s \n", e.Text)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.Visit("https://coding.imooc.com/learn/list/351.html")
}
