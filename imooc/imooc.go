package imooc

import (
	"fmt"
	"imooc_downloader/config"
	"imooc_downloader/tools"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
)

const (
	PRE_LOGIN_URL   = "https://www.imooc.com/passport/user/prelogin"
	SHOW_VERIFY_URL = "https://www.imooc.com/passport/user/loginverifyshow"
	LOGIN_URL       = "https://www.imooc.com/passport/user/login"
	VERIFY_CODE_URL = "https://www.imooc.com/passport/user/verifycode"
)

var referer = "https://www.imooc.com"

type UserManger struct {
	Username string
	Password string
}

var AuthCookies []*http.Cookie

// Gain access
func assignCookies(ssourl string) error {
	client := resty.New()
	url := tools.JoinAccessURL(ssourl)
	resp, err := client.R().Get(url)
	if err != nil {
		return fmt.Errorf("Gain access Get reqeuset failed. error: %v", err)
	}

	AuthCookies = resp.Cookies()
	return nil
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

	// 需要进行验证码图片的识别
	// if vermsg.Status == 10001 {
	// 	text := ocrByURL()
	// 	fmt.Printf("%v\n", text)
	// }

	form, err := u.createLoginForm(premsg)
	if err != nil {
		return err
	}

	subBody, err := submitForm(form)
	if err != nil {
		return err
	}
	lgmsg := new(LoginResponse)
	tools.Parser(subBody, lgmsg)
	if lgmsg.Status != 10001 {
		return fmt.Errorf(lgmsg.Msg)
	}

	urls, ok := lgmsg.Data["url"]
	if !ok {
		return fmt.Errorf("Field url does not exist.")
	}
	val, ok := urls.([]interface{})
	if !ok {
		return fmt.Errorf("Assertion url failed.")
	}

	url, ok := val[0].(string)
	if !ok {
		return fmt.Errorf("Assertion url[0] failed.")
	}

	// 获取 cookies
	assignCookies(url)

	return nil
}

func (u *UserManger) createLoginForm(premsg *PreLoginResponse) (map[string]string, error) {
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
		"browser_key": config.BROWSER_KEY,
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

func submitForm(data map[string]string) ([]byte, error) {
	client := resty.New()
	resp, err := client.R().
		SetHeaders(Headers).
		SetFormData(data).
		Post(LOGIN_URL)

	if err != nil {
		fmt.Printf("failed error: %v\n", err)
		return nil, err
	}

	return resp.Body(), nil
}

func handleDecryptPw(str string) ([]byte, error) {
	client := resty.New()
	resp, err := client.R().SetHeader("Content-type", "application/json").SetBody(map[string]interface{}{"pw": str}).Post(config.ENCRYPT_URL)
	if err != nil {
		return nil, fmt.Errorf("request decrypt password failed. error: %v\n", err)
	}
	return resp.Body(), nil
}

// func ocrByURL() string {
// 	text := ocr.ParserImgByURL(VERIFY_CODE_URL + "?t=" + strconv.Itoa(int(time.Now().UnixNano())/1e6))
// 	return text
// }
