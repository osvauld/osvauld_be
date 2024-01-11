package db

import (
	"context"
	"database/sql"
	dto "osvauld/dtos"

	"github.com/google/uuid"
)

type Store interface {
	Querier
	AddCredentialTransaction(context.Context, AddCredentialTransactionParams) (uuid.UUID, error)
	CreateGroupAndAddManager(context.Context, CreateGroupParams) (uuid.UUID, error)
	AddMemberToGroupTransaction(context.Context, AddMemberToGroupTransactionParams) error
	ShareCredentialWithUserTransaction(context.Context, dto.CredentialEncryptedFieldsForUserDto) error
	ShareCredentialWithGroupTransaction(context.Context, dto.CredentialEncryptedFieldsForGroupDto) error
	ShareMultipleCredentialsWithMultipleUsersTransaction(context.Context, []dto.CredentialEncryptedFieldsForUserDto) error
	ShareMultipleCredentialsWithMultipleGroupsTransaction(context.Context, []dto.CredentialEncryptedFieldsForGroupDto) error
	ShareFolderWithUsersTransaction(context.Context, uuid.UUID, []dto.CredentialsForUsersPayload) error
}

type SQLStore struct {
	Conn *sql.DB
	*Queries
}

func NewStore(conn *sql.DB) Store {
	return &SQLStore{
		Conn:    conn,
		Queries: New(conn),
	}
}
