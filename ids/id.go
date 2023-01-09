package ids

import (
	"encoding/hex"
	"errors"
	"strings"

	"github.com/thanhfphan/blockchain/utils/base58"
	"github.com/thanhfphan/blockchain/utils/hashing"
)

var (
	IDEmpty = ID{}
)

type ID [32]byte

func ToID(bytes []byte) (ID, error) {
	return hashing.ToHash256(bytes)
}

func IDFromString(idStr string) (ID, error) {
	bytes, err := base58.Decode(idStr)
	if err != nil {
		return IDEmpty, err
	}

	return ToID(bytes)
}

func IDFromPrefixedString(idStr, prefix string) (ID, error) {
	if !strings.HasPrefix(idStr, prefix) {
		return IDEmpty, errors.New("prefix is not correct")
	}
	return IDFromString(strings.TrimPrefix(idStr, prefix))
}

func (id ID) Bytes() []byte {
	return id[:]
}

func (id ID) Hex() string {
	return hex.EncodeToString(id.Bytes())
}

func (id ID) String() string {
	str, _ := base58.Encode(id.Bytes())
	return str
}

func (id ID) PrefixedString(prefix string) string {
	return prefix + id.String()
}
