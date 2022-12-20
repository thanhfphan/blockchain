package base58

import (
	"github.com/mr-tron/base58/base58"
)

func Encode(bytes []byte) (string, error) {
	return base58.Encode(bytes), nil
}

func Decode(str string) ([]byte, error) {
	return base58.Decode(str)
}
