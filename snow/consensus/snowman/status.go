package snowman

import (
	"errors"
)

type Status uint32

var errUnknowStatus = errors.New("unknow status")

const (
	Unknown Status = iota
	Processing
	Rejected
	Accepted
)

func (s Status) Decided() bool {
	switch s {
	case Rejected, Accepted:
		return true
	default:
		return false
	}
}

func (s Status) Valid() error {
	switch s {
	case Unknown, Processing, Rejected, Accepted:
		return nil
	default:
		return errUnknowStatus
	}
}

func (s Status) String() string {
	switch s {
	case Unknown:
		return "Unknown"
	case Processing:
		return "Processing"
	case Rejected:
		return "Rejected"
	case Accepted:
		return "Accepted"
	default:
		return "Invalid status"
	}
}
