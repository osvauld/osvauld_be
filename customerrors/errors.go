package customerrors

type UserNotAuthenticatedError struct {
	Message string
}

func (e *UserNotAuthenticatedError) Error() string {
	return e.Message
}
