package dl

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"imooc_downloader/common"
	"imooc_downloader/imooc"
	"imooc_downloader/tools"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/gosuri/uiprogress"
	"github.com/grafov/m3u8"
)

/*
	传入 m3u8 文件路径
	解析 seg 并下载, 然后用相应的 key 和 iv 来解码
	保存为 .ts 在本地
	用 ffmpeg 合并 .ts

*/

type Enginer struct {
	course    string
	dstFolder string
}

var tempFolderName = "course"

func NewEnginer(course string) *Enginer {
	return &Enginer{course: course, dstFolder: "./"}
}

func (en *Enginer) Download(chapter, name string, mediapl *m3u8.MediaPlaylist) {

	// c := http.Client{
	// 	Timeout: time.Duration(30) * time.Second,
	// }

	course, _ := filepath.Abs(filepath.Join(en.dstFolder, tempFolderName))
	segTemp := filepath.Join(course, chapter, name+"-"+time.Now().Format("20060102"))
	os.MkdirAll(segTemp, os.ModePerm)

	up := uiprogress.New()
	up.Start()
	bar := up.AddBar(int(mediapl.Count()))
	bar.AppendCompleted()

	for index, seg := range mediapl.Segments {
		if seg == nil {
			break
		}

		// resp, err := c.Get(seg.URI)
		// if err != nil {
		// 	return
		// }
		// if resp.StatusCode != 200 {
		// 	return
		// }

		// byteslice, _ := ioutil.ReadAll(resp.Body)
		// resp.Body.Close()
		d := Download{
			Url:           seg.URI,
			TotalSections: 5,
		}
		byteslice, err := d.Do()
		if err != nil {
			return
		}

		tag := seg.Custom[common.KEY_CONTENT_TAG]
		tagBytes, err := base64.StdEncoding.DecodeString(tag.String())

		iv := int64(mediapl.SeqNo + uint64(index))
		ivHex := fmt.Sprintf("%016x", iv)
		hx, _ := hex.DecodeString(ivHex)
		decryptedBytes, err := tools.Aes128Decrypt(byteslice, tagBytes, hx)

		tsFile := filepath.Join(segTemp, strconv.Itoa(index)+".ts")

		ioutil.WriteFile(tsFile, decryptedBytes, os.ModePerm)
		bar.Incr()
	}
	up.Stop()
	mergeRename(segTemp)
}

// 直接保存 markdown
func (en *Enginer) SaveAsFile(chapter, name, url string) {
	course, _ := filepath.Abs(filepath.Join(en.dstFolder, tempFolderName))
	segTemp := filepath.Join(course, chapter, name)
	os.MkdirAll(filepath.Join(course, chapter), os.ModePerm)

	client := resty.New()
	client.SetCookies(imooc.AuthCookies)

	resp, _ := client.R().SetHeaders(imooc.Headers).Get(url)

	docInfo := new(struct {
		Data map[string]interface{} `json:"data"`
	})

	tools.Parser(resp.Body(), docInfo)

	info := docInfo.Data["media_info"].(map[string]interface{})
	byteKey := []byte(fmt.Sprintf("%v", info["desc_md"]))
	ioutil.WriteFile(segTemp+".md", byteKey, os.ModePerm)
}

func mergeRename(src string) {

	if _, err := exec.LookPath("ffmpeg"); err == nil {
		files, err := ioutil.ReadDir(src)
		if err != nil {
			fmt.Printf("ReadDir failed. %v\n", err)
		}

		var line []byte = []byte{}
		fileCount := len(files)

		for index := range files {
			if index == fileCount-1 {
				line = append(line, []byte(src+"\\"+strconv.Itoa(index)+".ts")...)
			} else {
				line = append(line, []byte(src+"\\"+strconv.Itoa(index)+".ts"+"|")...)
			}
		}

		target := src + ".mp4"
		cmd := exec.Command("ffmpeg",
			"-i", "concat:"+string(line),
			"-y",
			"-acodec", "copy",
			"-vcodec", "copy",
			"-absf", "aac_adtstoasc",
			target,
		)
		err = cmd.Run()
		if err != nil {
			fmt.Printf("FFmpeg run command failed. %v\n", err)
		}

		err = os.RemoveAll(src)
		if err != nil {
			fmt.Printf("Remove temp failed. %v\n", err)
		}
	} else {
		fmt.Println("FFmpeg does not exist! skip merge")
	}
}

// 处理 windows 下文件路径过长
func prunePathLength(p string) string {

	// 通过 windows api 获取长度, 这里直接写个常量值
	if len(p) > 250 {

		return p
	}

	return p
}
