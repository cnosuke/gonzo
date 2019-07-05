package entity

import "time"

type CreatePresignedPostUrlParams struct {
	OriginalFilename string        `json:"filename"`
	Ttl              time.Duration `json:"ttl"`
}

type PresignedPostUrl struct {
	Url               string        `json:"url"`
	Ttl               time.Duration `json:"ttl"`
	Key               string        `json:"key"`
	Filename          string        `json:"filename"`
	ContentType       string        `json:"content_type"`
	PublicDownloadUrl string        `json:"public_download_url"`
}
