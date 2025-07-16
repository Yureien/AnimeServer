package utils

import (
	"github.com/zorchenhimer/go-ed2k"
)

func HashED2K(data []byte) []byte {
	hashed := ed2k.New()
	hashed.Write(data)
	return hashed.Sum(nil)
}
