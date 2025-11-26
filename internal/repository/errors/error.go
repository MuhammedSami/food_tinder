package errors

type NotFound struct {
	error
	err string
}

func NewNotFoundError(err string) *NotFound {
	return &NotFound{
		err: err,
	}
}
