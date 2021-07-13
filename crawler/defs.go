package crawler

import (
	"bytes"
)

type Lesson struct {
	title   string
	url     string
	chapter string
	isWork  bool
}

type msg struct {
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
