package api

type manager interface {
	CreateSession() (string, error)
}
