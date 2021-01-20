package tools

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"fmt"
	"imooc_downloader/config"
	"io"
	"net/http"
	"strconv"
	"time"
)

func Get(u string, params map[string]string) (io.ReadCloser, error) {
	c := http.Client{Timeout: time.Second * time.Duration(30)}
	request, err := http.NewRequest(http.MethodGet, u, nil)
	values := request.URL.Query()
	for k, v := range params {
		values.Add(k, v)
	}
	request.URL.RawQuery = values.Encode()

	for k, v := range config.FakeHeaders {
		request.Header.Add(k, v)
	}

	resp, err := c.Do(request)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http error: status code %d", resp.StatusCode)
	}

	return resp.Body, nil
}

func GetSegment(url string) ([]byte, error) {
	_, err := Get(url, nil)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func DecryptResponseInfo() {
}

func Parser(body []byte, v interface{}) error {
	decode := json.NewDecoder(bytes.NewBuffer(body))
	return decode.Decode(v)
}

func Aes128Decrypt(crypted, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	if len(iv) == 0 {
		iv = key
	}
	blockMode := cipher.NewCBCDecrypter(block, iv[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = Pkcs5UnPadding(origData)
	return origData, nil
}

func Pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

func JoinAccessURL(ssourl string) string {
	return ssourl + "&callback=jQuery19109012039733904491_" + strconv.Itoa(int(time.Now().UnixNano())/1e6) + "&_=" + strconv.Itoa(int(time.Now().UnixNano())/1e6+2)
}
