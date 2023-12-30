package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Store interface {
	Querier
	AddCredentialTransaction(context.Context, AddCredentialTransactionParams) (uuid.UUID, error)
	ShareCredentialTransaction(context.Context, ShareCredentialParams) error
	CreateGroupAndAddManager(context.Context, CreateGroupParams) (uuid.UUID, error)
	AddMemberToGroupTransaction(context.Context, AddMemberToGroupTransactionParams) error
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
