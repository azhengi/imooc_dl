package main

import (
	"errors"
	"fmt"
	"imooc_downloader/crawler"
	"imooc_downloader/execEnv"
	"imooc_downloader/imooc"
	"regexp"
	"strings"

	"github.com/manifoldco/promptui"
)

var emailReStr string = `^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`
var phoneReStr string = `^1(3\d|4[5-8]|5[0-35-9]|6[567]|7[01345-8]|8\d|9[025-9])\d{8}$`
var jsCodePath string = "./scripts/index.js"

var accountValidate = func(input string) error {
	trimStr := strings.Trim(input, " \n")
	emailMatched, _ := regexp.MatchString(emailReStr, trimStr)
	if emailMatched {
		return nil
	}

	phoneMatched, _ := regexp.MatchString(phoneReStr, trimStr)
	if phoneMatched {
		return nil
	}

	return errors.New("Invalid imooc account")
}

var pwdValidate = func(input string) error {
	return nil
}

func main() {

	execEnv.NewJsRuntime(jsCodePath)

	accountPrompt := promptui.Prompt{
		Label:       "账号",
		Validate:    accountValidate,
		HideEntered: true,
	}

	acc, err := accountPrompt.Run()
	if err != nil {
		fmt.Printf("account Prompt failed %v\n", err)
		return
	}

	pwdPrompt := promptui.Prompt{
		Label:       "密码",
		Validate:    pwdValidate,
		Mask:        '*',
		HideEntered: true,
	}
	pwd, err := pwdPrompt.Run()
	if err != nil {
		fmt.Printf("password Prompt failed %v\n", err)
		return
	}

	coursePrompt := promptui.Prompt{
		Label:   "课程链接( 暂不支持免费课程下载 )",
		Default: "https://coding.imooc.com/learn/list/99.html",
	}

	course, err := coursePrompt.Run()
	if err != nil {
		fmt.Printf("course Prompt failed %v\n", err)
		return
	}

	err = imooc.ParserCookieFile("authcookie")

	if err != nil {
		fmt.Printf("cookie login failed. error: %v\n", err)

		um := new(imooc.UserManger)
		um.Username = acc // phone or email to login
		um.Password = pwd
		err = um.DoLogin()
		if err != nil {
			fmt.Printf("run DoLogin failed. error: %v\n", err)
			return
		}
	}

	// do crawl
	crawler.StarColly(course)
}
