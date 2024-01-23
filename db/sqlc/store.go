package db

import (
	"context"
	"database/sql"
	dto "osvauld/dtos"

	"github.com/google/uuid"
)

type Store interface {
	Querier
	AddCredentialTransaction(context.Context, dto.AddCredentialDto, uuid.UUID) (uuid.UUID, error)
	CreateGroupAndAddManager(context.Context, CreateGroupParams) (uuid.UUID, error)
	AddMemberToGroupTransaction(context.Context, AddMemberToGroupTransactionParams) error
	ShareCredentialWithUserTransaction(context.Context, dto.CredentialFieldsForUserDto) error
	ShareCredentialWithGroupTransaction(context.Context, uuid.UUID, string, []dto.GroupCredentialPayload) error
	ShareMultipleCredentialsWithMultipleUsersTransaction(context.Context, []dto.CredentialFieldsForUserDto) error
	ShareMultipleCredentialsWithMultipleGroupsTransaction(context.Context, []dto.CredentialEncryptedFieldsForGroupDto) error
	ShareFolderWithUserTransaction(context.Context, uuid.UUID, dto.CredentialsForUsersPayload) error
	ShareFolderWithGroupTransaction(context.Context, uuid.UUID, dto.CredentialsForGroupsPayload) error
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
