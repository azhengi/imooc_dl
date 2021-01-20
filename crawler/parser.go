package crawler

import (
	"bytes"
	"imooc_downloader/config"
	"imooc_downloader/imooc"
	"imooc_downloader/tools"

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
func m3u8Parser(url string) *m3u8.MediaPlaylist {

	content := decryptURL(url, "")
	p, listType, _ := m3u8.DecodeFrom(bytes.NewReader(content), true)

	for listType == m3u8.MASTER {
		masterpl := p.(*m3u8.MasterPlaylist)
		variant := getMaxOfSlice(masterpl.Variants)
		content := decryptURL(variant.URI, "")
		p, listType, _ = m3u8.DecodeFrom(bytes.NewReader(content), true)
	}

	mediapl := p.(*m3u8.MediaPlaylist)
	if mediapl.Key != nil {
		// keyByte := decryptURL(mediapl.Key.URI, "1")
		// strKey := string(keyByte)
		/* for _, v := range strKey {
			fmt.Printf("%+v", string(v))
		} */
		// fmt.Printf("%+v\n ", string(strKey))
		// fmt.Printf("%+v\n", mediapl)
		// normalName := strings.Replace(name, ":", "_", -1)
		// ioutil.WriteFile(normalName+".m3u8", mediapl.Encode().Bytes(), os.ModePerm)
		return mediapl
	}

	return nil
}

func decryptURL(url, e string) []byte {
	client := resty.New()
	client.SetCookies(imooc.AuthCookies)

	resp, _ := client.R().SetHeaders(imooc.Headers).Get(url)
	pl := new(m3uResponse)
	tools.Parser(resp.Body(), pl)
	info := pl.Data["info"]
	resp, _ = client.R().SetHeader("Content-type", "application/json").SetBody(map[string]interface{}{"info": info, "e": e}).Post(config.DECRYPT_INFO_URL)
	return resp.Body()
}
