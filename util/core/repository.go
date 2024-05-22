package core

type ItemNotFoundException interface {
	Error() string
	IsItemNotFound() bool
}
