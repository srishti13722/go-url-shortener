package utils

import "github.com/segmentio/ksuid"

func GenerateShortCode() string{
	id := ksuid.New().String()[:6]
	return id
}