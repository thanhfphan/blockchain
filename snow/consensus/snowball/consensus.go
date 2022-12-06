package snowball

type Consensus interface {
	Initialize(initPrefernce int)

	Add(choiceID int)
	Preference() int
	Finalized() bool
}
