package customerrors

import (
	"fmt"

	"github.com/google/uuid"
)

type UserDoesNotHaveCredentialAccessError struct {
	UserID       uuid.UUID
	CredentialID uuid.UUID
}

func (e *UserDoesNotHaveCredentialAccessError) Error() string {
	return fmt.Sprintf("user %s does not have access to credential %s", e.UserID, e.CredentialID)
}

type UserDoesNotHaveFolderAccessError struct {
	UserID   uuid.UUID
	FolderID uuid.UUID
}

func (e *UserDoesNotHaveFolderAccessError) Error() string {
	return fmt.Sprintf("user %s does not have access to folder %s", e.UserID, e.FolderID)
}

type UserNotManagerOfCredentialError struct {
	UserID       uuid.UUID
	CredentialID uuid.UUID
}

func (e *UserNotManagerOfCredentialError) Error() string {
	return fmt.Sprintf("user %s does not have manager access for credential %s", e.UserID, e.CredentialID)
}

type UserNotManagerOfFolderError struct {
	UserID   uuid.UUID
	FolderID uuid.UUID
}

func (e *UserNotManagerOfFolderError) Error() string {
	return fmt.Sprintf("user %s does not have manager access for folder %s", e.UserID, e.FolderID)
}

type UserAlreadyMemberOfGroupError struct {
	UserID  uuid.UUID
	GroupID uuid.UUID
}

func (e *UserAlreadyMemberOfGroupError) Error() string {
	return fmt.Sprintf("user %s is already a member of group %s", e.UserID, e.GroupID)
}

type UserNotMemberOfGroupError struct {
	UserID  uuid.UUID
	GroupID uuid.UUID
}

func (e *UserNotMemberOfGroupError) Error() string {
	return fmt.Sprintf("user %s is not a member of group %s", e.UserID, e.GroupID)
}

type UserNotAdminOfGroupError struct {
	UserID  uuid.UUID
	GroupID uuid.UUID
}

func (e *UserNotAdminOfGroupError) Error() string {
	return fmt.Sprintf("user %s is not a admin of group %s", e.UserID, e.GroupID)
}
