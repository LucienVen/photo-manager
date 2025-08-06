package objects

type PhotoRecord struct {
	Filename   string   `json:"filename"`
	URL        string   `json:"url"`
	ThumbURL   string   `json:"thumb_url"`
	Path       string   `json:"path"`
	UploadedAt int64    `json:"uploaded_at"`
	Tags       []string `json:"tags"`
	Desc       string   `json:"desc"`
	SizeKB     int      `json:"size_kb"`
	Width      int      `json:"width"`
	Height     int      `json:"height"`
	Hash       string   `json:"hash"`
}

type PicGoResult struct {
	FileName string `json:"fileName"`
	URL      string `json:"url"`
}
