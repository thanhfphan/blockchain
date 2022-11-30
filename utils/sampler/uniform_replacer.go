package sampler

import (
	"errors"
	"math"
)

type uniformReplacer struct {
	rng        rng
	seededRNG  rng
	length     uint64
	drawn      map[uint64]uint64
	drawsCount uint64
}

func (u *uniformReplacer) Initialize(length uint64) error {
	if length > math.MaxInt64 {
		return errors.New("out of length Int64")
	}

	u.rng = globalRNG
	u.seededRNG = newRNG()
	u.length = length
	u.drawn = map[uint64]uint64{}
	u.drawsCount = 0

	return nil
}

func (u *uniformReplacer) Sample(count int) ([]uint64, error) {
	u.Reset()

	results := make([]uint64, count)
	for i := 0; i < count; i++ {
		ret, err := u.Next()
		if err != nil {
			return nil, err
		}
		results[i] = ret
	}
	return results, nil
}

func (u *uniformReplacer) Seed(seed int64) {
	u.rng = u.seededRNG
	u.rng.Seed(seed)
}

func (u *uniformReplacer) ClearSeed() {
	u.rng = globalRNG
}

func (u *uniformReplacer) Reset() {
	for k := range u.drawn {
		delete(u.drawn, k)
	}
	u.drawsCount = 0
}

func (u *uniformReplacer) Next() (uint64, error) {
	if u.drawsCount >= u.length {
		return 0, errors.New("out of length - unitformReplacer.Next")
	}

	draw := uint64(u.rng.Int63n(int64(u.length-u.drawsCount))) + u.drawsCount
	ret := draw
	if val, ok := u.drawn[draw]; ok {
		ret = val
	}
	u.drawn[draw] = u.drawsCount
	u.drawsCount++

	return ret, nil
}
