package crawler

import (
	"bytes"
)

type Lesson struct {
	title   string
	m3u8    string
	chapter string
}

type decryptMsg struct {
	Data map[string]interface{} `json:"data"`
}

type imoocKey string

func (ik imoocKey) TagName() string {
	return ""
}

func (ik imoocKey) Encode() *bytes.Buffer {
	return nil
}

func (ik imoocKey) String() string {
	return string(ik)
}
