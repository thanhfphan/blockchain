package sampler

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUniformReplacer_OverSample(t *testing.T) {
	r := require.New(t)
	s := NewUniform()

	err := s.Initialize(2)
	r.NoError(err)

	_, err = s.Sample(3)
	r.Error(err, "should returned an error")
}

func TestUniformReplacer_Singleton(t *testing.T) {
	r := require.New(t)
	s := NewUniform()

	err := s.Initialize(1)
	r.NoError(err)

	result, err := s.Sample(1)
	r.Equal([]uint64{0}, result, "should return the only element")
}

func TestUniformReplacer_DistributionTest(t *testing.T) {
	r := require.New(t)
	s := NewUniform()

	err := s.Initialize(3)
	r.NoError(err)

	result, err := s.Sample(3)
	sort.Slice(result, func(i, j int) bool {
		return result[i] < result[j]
	})
	r.Equal([]uint64{0, 1, 2}, result, "should return all elements")
}
