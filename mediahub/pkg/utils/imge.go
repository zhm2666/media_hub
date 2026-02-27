package utils

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
)

func IsImage(r io.Reader) bool {
	_, _, err := image.Decode(r)
	if err != nil {
		return false
	}
	return true
}
