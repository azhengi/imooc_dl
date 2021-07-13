package crawler

import (
	"fmt"
	"imooc_downloader/dl"
	"imooc_downloader/imooc"
	"log"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/grafov/m3u8"

	"github.com/gocolly/colly"
)

var lessonRe = regexp.MustCompile(`\/lesson\/(\d+)\.html#mid=(\d+)`)
var titleRe = regexp.MustCompile(`(\n| )+`)
var durRe = regexp.MustCompile(`(?:(\d{1,2}):)?(\d{1,2}):(\d{2})`)

var courseTitle string = ""

func StarColly(url string) {

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.113 Safari/537.36"),
	)

	ch := make(chan Lesson)
	ownTitleCh := make(chan struct{})
	wg := &sync.WaitGroup{}

	c.SetCookies("https://www.imooc.com", imooc.AuthCookies)

	c.OnHTML("h2[class=course-title]", func(e *colly.HTMLElement) {
		courseTitle = strings.Trim(e.Text, " \n")
		ownTitleCh <- struct{}{}
		log.Println(e.Text)
	})

	c.OnHTML("div[class=list-item]", func(e *colly.HTMLElement) {
		var chapter string

		e.ForEach("h3", func(i int, e *colly.HTMLElement) {
			chapter = strings.Trim(e.Text, " \n")
			log.Printf("%v\n", chapter)
		})

		e.ForEach("ul a", func(i int, e *colly.HTMLElement) {
			href := e.Attr("href")
			matches := lessonRe.FindAllStringSubmatch(href, -1)
			match := matches[0]
			cid := match[1]
			mid := match[2]

			var url string
			var isWork bool

			if isWork = !durRe.MatchString(e.Text); isWork {
				url = joinWorkDocURL(e.Request.URL, mid, cid)
			} else {
				url = joinM3u8H5URL(e.Request.URL, mid, cid)
			}

			title := titleRe.ReplaceAllString(e.Text, " ")
			normalName := strings.Replace(title, ":", "_", -1)
			lesson := Lesson{chapter: chapter, title: normalName, url: url, isWork: isWork}
			ch <- lesson
			wg.Add(1)

			log.Printf("  %v\n", title)
		})
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	c.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	go requestLesson(ownTitleCh, ch, wg)

	c.Visit(url)

	wg.Wait()
}

func requestLesson(ownTitleCh chan struct{}, dataChan chan Lesson, wg *sync.WaitGroup) {
	<-ownTitleCh
	en := dl.NewEnginer(courseTitle)

	for l := range dataChan {

		if l.isWork {
			en.SaveAsFile(l.chapter, l.title, l.url)
		} else {
			mediapl := m3u8Parser(l.url)
			en.Download(l.chapter, l.title, mediapl)
		}

		wg.Done()
	}
}

func joinM3u8H5URL(URL *url.URL, mid, cid string) string {
	hostname := URL.Host
	scheme := URL.Scheme
	return scheme + "://" + hostname + "/lesson/m3u8h5?mid=" + mid + "&cid=" + cid + "&ssl=1&cdn=aliyun1"
}

func joinWorkDocURL(URL *url.URL, mid, cid string) string {
	hostname := URL.Host
	scheme := URL.Scheme
	return scheme + "://" + hostname + "/lesson/ajaxmediainfo?mid=" + mid + "&cid=" + cid
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
