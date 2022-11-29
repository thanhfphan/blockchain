package handler

var _ Handler = (*handler)(nil)

type Handler interface {
}

type handler struct {
}

func New() (Handler, error) {
	h := &handler{}

	return h, nil
}
