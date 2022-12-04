package snowman

type Consensus interface {
	Initialize(initPreference int)
	Add(newChoice int)
	Preference() int
	// TODO: RecordPoll
	// TODO: RecordUnsuccessPoll
	Finalized() bool
}
