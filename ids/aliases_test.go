package ids

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAlias_E2E(t *testing.T) {
	r := require.New(t)
	a := NewAliaser()

	alias := "for the boys"
	testID := ID{33}

	notExistID, err := a.Lookup(alias)
	r.Error(err)
	r.Equal(IDEmpty, notExistID)

	aliases, err := a.Aliases(testID)
	r.NoError(err)
	r.Equal(0, len(aliases))

	err = a.Alias(testID, alias)
	r.NoError(err, "no error because ID is not exist")
	err = a.Alias(testID, alias)
	r.Error(err, "should throw error cause key already existed")

	aliases, err = a.Aliases(testID)
	r.NoError(err)
	r.Equal(1, len(aliases))

	existID, err := a.Lookup(alias)
	r.NoError(err)
	r.Equal(testID.String(), existID.String())

	a.RemoveAlias(existID)

	notExistID, err = a.Lookup(alias)
	r.Error(err)
	r.Equal(IDEmpty, notExistID, "not found because we just removed")

	aliases, err = a.Aliases(testID)
	r.NoError(err)
	r.Equal(0, len(aliases), "not found because we just removed")
}
