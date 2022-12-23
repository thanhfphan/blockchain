package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_AtomicInterface(t *testing.T) {
	ai := NewAtomicInterface(nil)
	require.Nil(t, ai.GetValue())

	ai.SetValue("hello")
	require.Equal(t, "hello", ai.GetValue().(string))
}
