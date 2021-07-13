package crawler

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"imooc_downloader/common"
	"imooc_downloader/config"
	"imooc_downloader/execEnv"
	"imooc_downloader/imooc"
	"imooc_downloader/tools"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/grafov/m3u8"
)

/*

1. 先通过 url 得到 playlist info 再进行解密得到 playlist m3u8
2. 提取 playlist m3u8 得到指定清晰度的 resolution url
3. 请求 resolution url 得到最终视频文件的 video info 并且进行解密得到 video m3u8
4. 提取 video m3u8 中的 KEY URL.  请求 KEY URL 对响应再进行解密得到实际 KEY
5. 提取 video m3u8 得到各个 segment 的 .ts 路径
6. 依次下载 .ts 并用实际 KEY 来解密. 得到最终可播放的 .ts 文件


从 3 开始, 可以省略剩余步骤, 直接将 video m3u8 通过 N_m3u8DL-CLI 运行.
https://github.com/nilaoda/N_m3u8DL-CLI

*/

var valueRe = regexp.MustCompile(`:(\d+)`)

func m3u8Parser(url string) *m3u8.MediaPlaylist {

	content := decryptURLOfLocal(url, "")
	p, listType, err := m3u8.DecodeFrom(bytes.NewReader(content), true)
	if err != nil {
		fmt.Printf("m3u8 decodeFrom failed %v\n", err)
		return nil
	}

	for listType == m3u8.MASTER {
		masterpl, _ := p.(*m3u8.MasterPlaylist)
		variant := getMaxOfSlice(masterpl.Variants)
		// 获取指定清晰度 m3u8 url . variant.URI
		content := decryptURLOfLocal(variant.URI, "")
		p, listType, _ = m3u8.DecodeFrom(bytes.NewReader(content), true)
	}

	mediapl, ok := p.(*m3u8.MediaPlaylist)
	if ok {
		var segDecodeKey string

		for _, seg := range mediapl.Segments {
			if seg == nil {
				continue
			}

			if seg.Key != nil {
				segDecodeKey = keyParserLocalStr(seg.Key.URI)
			}
			if seg.Custom == nil {
				seg.Custom = make(map[string]m3u8.CustomTag)
			}
			seg.Custom[common.KEY_CONTENT_TAG] = imoocKey(segDecodeKey)
		}
	} else {

		log.Printf("p: %+v, 断言失败 mediapl %+v:\n", p, mediapl)
	}

	return mediapl
}

func decryptURL(url, e string) []byte {
	client := resty.New()
	client.SetCookies(imooc.AuthCookies)

	r := client.R().SetHeaders(imooc.Headers)
	resp, err := r.Get(url)
	if err != nil {
		fmt.Printf("client Get failed %v\n", err)
		return nil
	}
	pl := new(msg)
	tools.Parser(resp.Body(), pl)
	info := pl.Data["info"]
	resp, err = client.R().SetHeader("Content-type", "application/json").SetBody(map[string]interface{}{"info": info, "e": e}).Post(config.DECRYPT_INFO_URL)
	if err != nil {
		fmt.Printf("client Post failed %v\n", err)
		return nil
	}
	return resp.Body()
}

func decryptURLOfLocal(url, e string) []byte {
	client := resty.New()
	client.SetCookies(imooc.AuthCookies)

	r := client.R().SetHeaders(imooc.Headers)
	resp, err := r.Get(url)
	if err != nil {
		fmt.Printf("client Get failed %v\n", err)
		return nil
	}
	pl := new(msg)
	tools.Parser(resp.Body(), pl)
	info := pl.Data["info"]
	infoStr := info.(string)
	jsResult := execEnv.JsRt.DecryptInfo(infoStr, e)
	return jsResult
}

func keyParser(uri string) string {
	content := decryptURL(uri, "1")
	str := string(content)
	matches := valueRe.FindAllStringSubmatch(str, -1)
	values := make([]byte, 0, 16)

	for _, v := range matches {
		val, _ := strconv.ParseInt(v[1], 10, 64)
		values = append(values, byte(val))
	}

	encodeString := base64.StdEncoding.EncodeToString(values)
	return encodeString
}

func keyParserLocalStr(uri string) string {
	content := decryptURLOfLocal(uri, "1")
	str := string(content)

	sl := strings.Split(str, ",")

	values := make([]byte, 0, 16)

	for _, v := range sl {
		i, _ := strconv.Atoi(v)
		values = append(values, byte(i))
	}

	encodeString := base64.StdEncoding.EncodeToString(values)
	return encodeString
}
