package ocr

import (
	"github.com/otiai10/gosseract/v2"
)

// ParserImgByURL get verify code
func ParserImgByURL(url string) string {
	client := gosseract.NewClient()
	defer client.Close()
	client.SetImage(url)
	text, _ := client.Text()
	return text
}
