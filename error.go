package app

// Error represents app internal errors.
type Error string

// Error returns the error message. Fulfills the error interface
func (e Error) Error() string { return string(e) }

// user errors
const (
	ErrEmailAlreadyUsed    = Error("email already in use")
	ErrWrongPasswordFormat = Error("wrong password format")
	ErrUserNotFound        = Error("not found")
	ErrWrongCredentials    = Error("wrong credentials")
)

// article errors
const (
	ErrArticleNotFound = Error("not found")
)
