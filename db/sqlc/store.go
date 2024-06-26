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
	CreateFolderTransaction(context.Context, CreateFolderTransactionParams) (dto.FolderDetails, error)
	CreateGroupAndAddManager(context.Context, dto.GroupDetails) (dto.GroupDetails, error)
	AddMembersToGroupTransaction(context.Context, AddMembersToGroupTransactionParams) error
	ShareCredentialsTransaction(context.Context, ShareCredentialTransactionParams) error
	EditCredentialTransaction(context.Context, dto.EditCredentialRequest, uuid.UUID) error
	RemoveFolderAccessForUsersTransactions(context.Context, RemoveFolderAccessForUsersParams) error
	RemoveFolderAccessForGroupsTransactions(context.Context, RemoveFolderAccessForGroupsParams) error
	EditFolderAccessForUserTransaction(context.Context, EditFolderAccessForUserParams) error
	EditFolderAccessForGroupTransaction(context.Context, EditFolderAccessForGroupParams) error
	RemoveMemberFromGroupTransaction(context.Context, RemoveMemberFromGroupTransactionParams) error
	AddCredentialFieldToEnvTxn(context.Context, []dto.CredentialEnvData) error
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
