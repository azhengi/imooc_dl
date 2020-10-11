package crawler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"

	"github.com/gocolly/colly"
)

// Gain access
func GainAccess(ssourl string) ([]*http.Cookie, error) {
	client := resty.New()
	resp, err := client.R().Get(ssourl + "&callback=jQuery19109012039733904491_" + strconv.Itoa(int(time.Now().UnixNano())/1e6) + "&_=" + strconv.Itoa(int(time.Now().UnixNano())/1e6+2))
	if err != nil {
		return nil, fmt.Errorf("Gain access Get reqeuset failed. error: %v", err)
	}

	cookies := resp.Cookies()
	return cookies, nil
}

func StarColly(ssourl string) {
	cookies, err := GainAccess(ssourl)
	if err != nil {
		return
	}

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.113 Safari/537.36"),
	)

	c.SetCookies("https://www.imooc.com", cookies)

	c.OnHTML("div[class=list-item]", func(e *colly.HTMLElement) {
		e.ForEach("h3", func(i int, e *colly.HTMLElement) {
			fmt.Printf("%v\n", strings.Trim(e.Text, " \n"))
		})

		e.ForEach("ul a", func(i int, e *colly.HTMLElement) {
			fmt.Printf("  %v\n", strings.Trim(e.Text, " \n"))
		})
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)
		// log.Println("response received Body", string(r.Body))
		// ioutil.WriteFile("./imooc_coding.html", r.Body, os.ModePerm)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.Visit("https://coding.imooc.com/learn/list/351.html")
}
