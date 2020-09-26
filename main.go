package main

import "imooc_downloader/imooc"

func main() {
	um := new(imooc.UserManger)
	um.Username = "" // phone or email to login
	um.Password = ""
	um.DoLogin()
}
