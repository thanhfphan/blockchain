package snowman

var (
	_ Consensus = (*Topological)(nil)
)

type Topological struct {
	pollNumber int
	head       int
	tail       int
	height     int
	blocks     map[int]*snowmanBlock
	kahnNodes  map[int]kahnNode
}

type kahnNode struct {
	degree int
}

func (t *Topological) Initialize() error {
	return nil
}

func (t *Topological) Add(newChoice int) {
}

func (t *Topological) Preference() int {
	return -1
}

// TODO: RecordPoll
// TODO: RecordUnsuccessPoll

func (t *Topological) Finalized() bool {
	return false
}
