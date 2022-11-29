package handler

type IHandler interface {
}

type Handler struct {
}

func New() (IHandler, error) {
	h := &Handler{}

	return h, nil
}
