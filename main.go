package main

import (
	"fmt"
	"imooc_downloader/crawler"
	"imooc_downloader/imooc"
)

func main() {
	um := new(imooc.UserManger)
	um.Username = "" // phone or email to login
	um.Password = ""
	ssourl, err := um.DoLogin()
	if err != nil {
		fmt.Printf("run DoLogin failed. error: %v\n", err)
		return
	}

	crawler.StarColly(ssourl)
}
