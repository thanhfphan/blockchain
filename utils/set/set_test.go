package set

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thanhfphan/blockchain/ids"
)

func Test_E2E(t *testing.T) {
	r := require.New(t)
	id1 := ids.ID{1}
	id2 := ids.ID{2}
	id3 := ids.ID{3}

	s := New[ids.ID](2)

	r.Equal(0, s.Len())
	s.Add(id1)
	r.Equal(1, s.Len())
	r.Equal(true, s.Contains(id1))

	s.Add(id2)
	s.Add(id3)
	r.Equal(3, s.Len())
	r.Equal(true, s.Contains(id2))
	r.Equal(true, s.Contains(id3))

	s.Clear()
	r.Equal(0, s.Len())
	r.Equal(false, s.Contains(id1))
	r.Equal(false, s.Contains(id2))
	r.Equal(false, s.Contains(id3))

	s.Add(id1)
	s.Add(id2)
	s.Add(id3)
	list := s.List()
	r.Equal(3, len(list))

}
