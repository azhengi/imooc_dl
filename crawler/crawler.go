package crawler

import (
	"log"

	"github.com/gocolly/colly"
)

func StarColly(ssourl string) {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.113 Safari/537.36"),
	)

	if err := c.Visit(ssourl); err != nil {
		log.Printf("visit ssourl. error: %s \n", err)
		return
	}

	c.OnHTML("h3", func(e *colly.HTMLElement) {
		log.Printf("section %s \n", e.Text)
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)
		// ioutil.WriteFile("./imooc_coding.html", r.Body, os.ModePerm)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.Visit("https://coding.imooc.com/learn/list/351.html")
}
