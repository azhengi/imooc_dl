package crawler

type Lesson struct {
	title   string
	m3u8    string
	chapter string
}

type playlistResponse struct {
	Result string                 `json:"result"`
	Data   map[string]interface{} `json:"data"`
	Msg    string                 `json:"msg"`
}

type m3uResponse struct {
	Data map[string]interface{} `json:"data"`
}
