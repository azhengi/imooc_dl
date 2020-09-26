package imooc

import (
	"fmt"
	"imooc_downloader/tools"
	"strconv"

	"github.com/go-resty/resty/v2"
)

const (
	PRE_LOGIN_URL   = "https://www.imooc.com/passport/user/prelogin"
	SHOW_VERIFY_URL = "https://www.imooc.com/passport/user/loginverifyshow"
	LOGIN_URL       = "https://www.imooc.com/passport/user/login"
	BROWSER_KEY     = "dd9eeccdd46ca5935707f07fef4ba2fb"
	DECRYPT_URL     = "http://34.80.19.136:58000/decrypt"
)

var referer = "https://www.imooc.com"

type UserManger struct {
	Username string
	Password string
}

func (u *UserManger) DoLogin() error {
	preBody, err := ready()
	if err != nil {
		return err
	}
	premsg := new(PreLoginResponse)
	tools.Parser(preBody, premsg)

	verBody, err := verify(u.Username)
	if err != nil {
		return err
	}
	vermsg := new(VerifyResponse)
	tools.Parser(verBody, vermsg)

	form, err := u.createLoginForm(premsg)
	if err != nil {
		return err
	}

	submitForm(form)

	return nil
}

func (u *UserManger) createLoginForm(premsg *PreLoginResponse) (map[string]string, error) {
	// 请求解密 code servertime password 并且 base64
	str := premsg.Code + "\t" + strconv.Itoa(premsg.ServerTime) + "\t" + u.Password
	pw, err := handleDecryptPw(str)
	if err != nil {
		return nil, err
	}

	// 再修改
	form := map[string]string{
		"username":    u.Username,
		"password":    string(pw),
		"verify":      "",
		"remember":    "1",
		"pwencode":    "1",
		"browser_key": BROWSER_KEY,
		"referer":     referer,
	}

	return form, nil
}

func ready() ([]byte, error) {
	client := resty.New()
	resp, err := client.R().
		SetHeaders(Headers).
		Post(PRE_LOGIN_URL)
	if err != nil {
		return nil, fmt.Errorf("failed error: %v\n", err)
	}

	return resp.Body(), nil
}

func verify(user string) ([]byte, error) {
	client := resty.New()
	resp, err := client.R().
		SetHeaders(Headers).
		Get(SHOW_VERIFY_URL + "?username=" + user)
	if err != nil {
		return nil, fmt.Errorf("failed error: %v\n", err)
	}

	return resp.Body(), nil
}

func submitForm(data map[string]string) {
	client := resty.New()
	resp, err := client.R().
		SetHeaders(Headers).
		SetFormData(data).
		Post(LOGIN_URL)

	if err != nil {
		fmt.Printf("failed error: %v\n", err)
		return
	}

	fmt.Printf("%v \n", resp)
}

func handleDecryptPw(str string) ([]byte, error) {
	client := resty.New()
	resp, err := client.R().SetHeader("Content-type", "application/json").SetBody(map[string]interface{}{"pw": str}).Post(DECRYPT_URL)
	if err != nil {
		return nil, fmt.Errorf("request decrypt password failed. error: %v", err)
	}
	return resp.Body(), nil
}
