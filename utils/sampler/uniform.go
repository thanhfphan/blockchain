package sampler

type Uniform interface {
	Initialize(sampleRrange uint64) error
	Sample(length int) ([]uint64, error)

	Seed(int64)
	ClearSeed()
	Reset()
	Next() (uint64, error)
}

func NewUniform() Uniform {
	return &uniformReplacer{}
}
