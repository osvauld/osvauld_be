package customerrors

type UserNotAuthenticatedError struct {
	Message string
}

func (e *UserNotAuthenticatedError) Error() string {
	return e.Message
}

type UserAlreadyMemberOfGroupError struct {
	Message string
}

func (e *UserAlreadyMemberOfGroupError) Error() string {
	return e.Message
}
