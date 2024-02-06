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
