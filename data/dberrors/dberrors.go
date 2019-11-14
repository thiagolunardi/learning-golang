package dberrors

// Error -
type Error string

func (e Error) Error() string { return string(e) }

// Errors -
const (
	ErrItemNotFound = Error("Item not found")
)
