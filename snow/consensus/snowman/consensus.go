package snowman

type Consensus interface {
	Initialize() error
	Add(newChoice int)
	Preference() int
	// TODO: RecordPoll
	// TODO: RecordUnsuccessPoll
	Finalized() bool
}
