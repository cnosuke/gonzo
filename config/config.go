package config

import (
	"github.com/jinzhu/configor"
)

type Config struct {
	AwsRegion           string `default:"ap-northeast-1" env:"AWS_REGION"`
	BucketName          string `required:"true" env:"BUCKET_NAME"`
	KeyPrefix           string `default:"" env:"KEY_PREFIX"`
	DefaultPresignedTTL int    `default:"60" env:"DEFAULT_PRESIGNED_TTL"`
	DownloadHost string `default:"" env:"DOWNLOAD_HOST"`
}

func NewConfig(path string) (*Config, error) {
	c := &Config{}

	if err := configor.Load(c, path); err != nil {
		return nil, err
	}

	return c, nil
}
