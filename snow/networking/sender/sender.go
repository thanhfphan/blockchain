package sender

var _ Sender = (*sender)(nil)

type Sender interface {
}

type sender struct {
}

func New() (Sender, error) {
	s := &sender{}

	return s, nil
}
