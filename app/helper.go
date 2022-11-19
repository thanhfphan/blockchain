package app

import "strconv"

func IntToHex(data int64) []byte {
	return []byte(strconv.FormatInt(data, 16))
}
