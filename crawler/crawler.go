package crawler

import (
	"fmt"
	"imooc_downloader/imooc"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/grafov/m3u8"

	"github.com/gocolly/colly"
)

var lessonRe = regexp.MustCompile(`\/lesson\/(\d+)\.html#mid=(\d+)`)

func StarColly(url string) {

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.113 Safari/537.36"),
	)

	ch := make(chan Lesson)

	c.SetCookies("https://www.imooc.com", imooc.AuthCookies)

	c.OnHTML("div[class=list-item]", func(e *colly.HTMLElement) {
		var chapter string

		e.ForEach("h3", func(i int, e *colly.HTMLElement) {
			chapter = strings.Trim(e.Text, " \n")
			fmt.Printf("%v\n", chapter)
		})

		e.ForEach("ul a", func(i int, e *colly.HTMLElement) {
			href := e.Attr("href")
			matches := lessonRe.FindAllStringSubmatch(href, -1)
			match := matches[0]
			cid := match[1]
			mid := match[2]

			m3u8URL := joinM3u8H5URL(e.Request.URL, mid, cid)

			title := strings.Trim(e.Text, " \n")
			lesson := Lesson{chapter: chapter, title: title, m3u8: m3u8URL}
			ch <- lesson

			fmt.Printf("  %v\n", title)
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

	go requestLesson(ch)

	c.Visit(url)
}

func requestLesson(dataChan chan Lesson) {
	for {
		func() {
			select {
			case lesson := <-dataChan:
				mediapl := m3u8Parser(lesson.m3u8)
				// os.MkdirAll(lesson.chapter, os.ModePerm)
				normalName := strings.Replace(lesson.title, ":", "_", -1)
				ioutil.WriteFile(normalName+".m3u8", mediapl.Encode().Bytes(), os.ModePerm)
			}
		}()
	}
}

func joinM3u8H5URL(URL *url.URL, mid, cid string) string {
	hostname := URL.Host
	scheme := URL.Scheme
	return scheme + "://" + hostname + "/lesson/m3u8h5?mid=" + mid + "&cid=" + cid + "&ssl=1&cdn=aliyun1"
}

func getMaxOfSlice(sl []*m3u8.Variant) *m3u8.Variant {
	if len(sl) > 0 {
		max := sl[0]
		for _, variant := range sl {
			if variant.Bandwidth > max.Bandwidth {
				fmt.Printf("Resolution: %+v, bandwidth: %+v \n", variant.Resolution, variant.Bandwidth)
				max = variant
			}
		}
		return max
	}
	return nil
}
