package imooc

import (
	"fmt"
	"imooc_downloader/tools"

	"github.com/go-resty/resty/v2"
)

const (
	PRE_LOGIN_URL   = "https://www.imooc.com/passport/user/prelogin"
	SHOW_VERIFY_URL = "https://www.imooc.com/passport/user/loginverifyshow"
	LOGIN_URL       = "https://www.imooc.com/passport/user/login"
	BROWSER_KEY     = "dd9eeccdd46ca5935707f07fef4ba2fb"
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

	preMsg := new(PreLoginResponse)
	tools.Parser(preBody, preMsg)

	fmt.Printf(
		`code: %v
pubkey: %v
server time: %v
status: %v`,
		preMsg.Code, preMsg.PubKey, preMsg.ServerTime, preMsg.Status)
	fmt.Println()

	verBody, err := verify(u.Username)
	if err != nil {
		return err
	}

	verMsg := new(VerifyResponse)
	tools.Parser(verBody, verMsg)

	fmt.Printf(
		`msg: %v
need verify code ?: %v`,
		verMsg.Msg, verMsg.Status == 10001)
	fmt.Println()

	signInForm := map[string]string{
		"username":    u.Username,
		"password":    "SFHIkqckHWxZb/qJp1tRZEUVkJJCJ3CBLxy18rYJ1GqF7MD+e8B7Gkwud1P//NhIncgrzBA6SOAHVm9A/B28Df+QfENVZNK1Z7Z+8yBtR+ceUioFWyJhxPHrBXDOSqG95KCDMHjbJdF+zC5z26/d5O1xWKX53M5+Jlgb+z2mjy4=",
		"verify":      "",
		"remember":    "1",
		"pwencode":    "1",
		"browser_key": BROWSER_KEY,
		"referer":     referer,
	}
	signIn(signInForm)

	return nil
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

func signIn(data map[string]string) {
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
