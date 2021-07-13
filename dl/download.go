package dl

// by https://github.com/usmanhalalit/go-download-manager/blob/master/main.go

import (
	"errors"
	"fmt"
	"imooc_downloader/config"
	"imooc_downloader/imooc"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

type Download struct {
	Url           string
	TargetPath    string
	TotalSections int
}

// Start the download
func (d Download) Do() ([]byte, error) {

	client := resty.New()
	client.RetryCount = 3
	resp, _ := client.R().SetHeaders(imooc.Headers).Head(d.Url)

	if resp.StatusCode() > 299 {
		return nil, errors.New(fmt.Sprintf("Can't process, response is %v", resp.StatusCode()))
	}

	size, err := strconv.Atoi(resp.Header().Get("Content-Length"))
	if err != nil {
		return nil, err
	}

	var sections = make([][2]int, d.TotalSections)
	eachSize := size / d.TotalSections

	// example: if file size is 100 bytes, our section should like:
	// [[0 10] [11 21] [22 32] [33 43] [44 54] [55 65] [66 76] [77 87] [88 98] [99 99]]
	for i := range sections {
		if i == 0 {
			// starting byte of first section
			sections[i][0] = 0
		} else {
			// starting byte of other sections
			sections[i][0] = sections[i-1][1] + 1
		}

		if i < d.TotalSections-1 {
			// ending byte of other sections
			sections[i][1] = sections[i][0] + eachSize
		} else {
			// ending byte of other sections
			sections[i][1] = size - 1
		}
	}

	var wg sync.WaitGroup
	// download each section concurrently
	for i, s := range sections {
		wg.Add(1)
		go func(i int, s [2]int) {
			defer wg.Done()
			err = d.downloadSection(i, s)
			if err != nil {
				panic(err)
			}
		}(i, s)
	}
	wg.Wait()

	byteSlice, _ := d.mergeFiles(sections)
	return byteSlice, nil
}

// Download a single section and save content to a tmp file
func (d Download) downloadSection(i int, c [2]int) error {

	client := resty.New()
	client.RetryCount = 3
	client.SetTimeout(time.Duration(30) * time.Second)

	r := client.R()
	mergeHeaders := map[string]string{}
	mergeHeaders["Range"] = fmt.Sprintf("bytes=%v-%v", c[0], c[1])
	for k, v := range imooc.Headers {
		mergeHeaders[k] = v
	}
	resp, _ := r.SetHeaders(mergeHeaders).Get(d.Url)

	if resp.StatusCode() > 299 {
		return errors.New(fmt.Sprintf("Can't process, response is %v", resp.StatusCode()))
	}

	body := resp.Body()

	err := ioutil.WriteFile(fmt.Sprintf("section-%v.tmp", i), body, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// Get a new http request
func (d Download) getNewRequest(method string) (*http.Request, error) {
	r, err := http.NewRequest(
		method,
		d.Url,
		nil,
	)
	if err != nil {
		return nil, err
	}

	for k, v := range config.FakeHeaders {
		r.Header.Add(k, v)
	}
	return r, nil
}

// Merge tmp files to single file and delete tmp files
func (d Download) mergeFiles(sections [][2]int) ([]byte, error) {
	var sliceWrite = []byte{}
	for i := range sections {
		tmpFileName := fmt.Sprintf("section-%v.tmp", i)
		b, err := ioutil.ReadFile(tmpFileName)
		if err != nil {
			return nil, err
		}
		sliceWrite = append(sliceWrite, b...)
		err = os.Remove(tmpFileName)
		if err != nil {
			return nil, err
		}
	}

	return sliceWrite, nil
}
