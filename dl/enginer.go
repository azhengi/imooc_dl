package dl

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"imooc_downloader/common"
	"imooc_downloader/tools"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gosuri/uiprogress"
	"github.com/grafov/m3u8"
)

/*
	传入 m3u8 文件路径
	解析 seg 并下载, 然后用相应的 key 和 iv 来解码
	保存为 .ts 在本地
	用 ffmpeg 合并 .ts

*/

var dstFolder = "./dl_dst"

type Enginer struct {
	name    string
	mediapl *m3u8.MediaPlaylist
}

func NewEnginer(name string, mediapl *m3u8.MediaPlaylist) *Enginer {
	return &Enginer{name: name, mediapl: mediapl}
}

func (en *Enginer) Download() {

	c := http.Client{
		Timeout: time.Duration(30) * time.Second,
	}

	dstAbs, _ := filepath.Abs(dstFolder)
	segmentsFolder := filepath.Join(dstAbs, en.name, "segments")
	os.MkdirAll(segmentsFolder, os.ModePerm)

	uiprogress.Start()
	bar := uiprogress.AddBar(int(en.mediapl.Count()))
	bar.AppendCompleted()
	bar.PrependElapsed()
	pgs := 1

	for index, seg := range en.mediapl.Segments {
		if seg == nil {
			break
		}

		resp, err := c.Get(seg.URI)
		if err != nil {
			return
		}
		if resp.StatusCode != 200 {
			return
		}

		byteslice, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		tag := seg.Custom[common.KEY_CONTENT_TAG]
		tagBytes, err := base64.StdEncoding.DecodeString(tag.String())

		iv := int64(en.mediapl.SeqNo + uint64(index))
		ivHex := fmt.Sprintf("%016x", iv)
		hx, _ := hex.DecodeString(ivHex)
		decryptedBytes, err := tools.Aes128Decrypt(byteslice, tagBytes, hx)

		tsFile := filepath.Join(segmentsFolder, strconv.Itoa(index)+".ts")

		ioutil.WriteFile(tsFile, decryptedBytes, os.ModePerm)
		bar.Set(pgs)
		pgs++
	}
}
