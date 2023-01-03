package peer

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/thanhfphan/blockchain/utils/wrappers"
)

func writeMsgLen(msgLen uint32, maxMsgLen uint32) ([wrappers.IntLen]byte, error) {
	if msgLen > maxMsgLen {
		return [wrappers.IntLen]byte{}, errors.New("writeMsglen failed cause message length exceed the limit")
	}

	x := msgLen
	b := [wrappers.IntLen]byte{}
	binary.BigEndian.PutUint32(b[:], x)

	return b, nil
}

func readMsgLen(b []byte, maxMsgLen uint32) (uint32, error) {
	if len(b) != wrappers.IntLen {
		return 0, fmt.Errorf("readMsgLen only support 4-bytes - got %d bytes", len(b))
	}

	msgLen := binary.BigEndian.Uint32(b)
	if msgLen > maxMsgLen {
		return 0, fmt.Errorf("readMsgLen the message length %d exceeds the limit %d", msgLen, maxMsgLen)
	}
	return msgLen, nil
}
