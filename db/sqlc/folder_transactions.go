package db

import (
	"context"
	"database/sql"
	dto "osvauld/dtos"

	"github.com/google/uuid"
)

type CreateFolderTransactionParams struct {
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
	CreatedBy   uuid.UUID      `json:"createdBy"`
}

func (store *SQLStore) CreateFolderTransaction(ctx context.Context, args CreateFolderTransactionParams) (dto.FolderDetails, error) {

	var folderDetails dto.FolderDetails
	err := store.execTx(ctx, func(q *Queries) error {

		// Create folder record
		addFolderParams := AddFolderParams(args)
		folderData, err := q.AddFolder(ctx, addFolderParams)
		if err != nil {
			return err
		}

		// Add record to folder access table
		err = q.AddFolderAccess(ctx, AddFolderAccessParams{
			FolderID:   folderData.ID,
			UserID:     args.CreatedBy,
			AccessType: "owner",
		})
		if err != nil {
			return err
		}

		folderDetails.FolderID = folderData.ID
		folderDetails.CreatedAt = folderData.CreatedAt
		folderDetails.Name = args.Name
		folderDetails.Description = args.Description.String
		folderDetails.CreatedBy = args.CreatedBy

		return nil
	})

	return folderDetails, err
}
