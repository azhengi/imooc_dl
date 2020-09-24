package imooc

var Headers = map[string]string{
	"Accept":           "application/json, text/javascript, */*; q=0.01",
	"Content-Type":     "application/x-www-form-urlencoded; charset=UTF-8",
	"Referer":          "https://www.imooc.com/",
	"User-Agent":       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.113 Safari/537.36",
	"X-Requested-With": "XMLHttpRequest",
}

type PreLoginResponse struct {
	Code       string `json:"code"`
	PubKey     string `json:"pubkey"`
	ServerTime int    `json:"servertime"`
	Status     int    `json:"status"`
}

type VerifyResponse struct {
	Status int      `json:"status"`
	Msg    string   `json:"msg"`
	Data   []string `json:"data"`
}

type LoginResponse struct {
}
