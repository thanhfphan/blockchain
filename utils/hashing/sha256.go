package hashing

import (
	"crypto/sha256"
	"fmt"
)

const (
	Hash256Len = sha256.Size
)

type Hash256 = [Hash256Len]byte

func ToHash256(bytes []byte) (Hash256, error) {
	hash := Hash256{}

	if len(bytes) != Hash256Len {
		return hash, fmt.Errorf("expect 32 bytes but got %d", len(bytes))
	}

	copy(hash[:], bytes)
	return hash, nil
}

func ComputeHash256Array(buf []byte) Hash256 {
	return sha256.Sum256(buf)
}

func ComputeHash256(buf []byte) []byte {
	arr := ComputeHash256Array(buf)
	return arr[:]
}
