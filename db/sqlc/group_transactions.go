package db

import (
	"context"
	dto "osvauld/dtos"

	"github.com/google/uuid"
)

func (store *SQLStore) CreateGroupAndAddManager(ctx context.Context, args CreateGroupParams) (uuid.UUID, error) {

	var groupID uuid.UUID
	err := store.execTx(ctx, func(q *Queries) error {
		groupID, err := q.CreateGroup(ctx, CreateGroupParams{
			Name:      args.Name,
			CreatedBy: args.CreatedBy,
		})
		if err != nil {
			return err
		}

		err = q.AddGroupMemberRecord(ctx, AddGroupMemberRecordParams{
			GroupingID: groupID,
			UserID:     args.CreatedBy,
			AccessType: "manager",
		})
		if err != nil {
			return err
		}

		return nil
	})
	return groupID, err
}

type AddMemberToGroupTransactionParams struct {
	GroupID           uuid.UUID                          `json:"groupId"`
	UserID            uuid.UUID                          `json:"userId"`
	MemberRole        string                             `json:"memberRole"`
	UserEncryptedData []dto.CredentialEncryptedFieldsDto `json:"encryptedFields"`
}

func (store *SQLStore) AddMemberToGroupTransaction(ctx context.Context, args AddMemberToGroupTransactionParams) error {
	err := store.execTx(ctx, func(q *Queries) error {

		// Add record to grouping table
		err := q.AddGroupMemberRecord(ctx, AddGroupMemberRecordParams{
			GroupingID: args.GroupID,
			UserID:     args.UserID,
			AccessType: args.MemberRole,
		})
		if err != nil {
			return err
		}

		// Add Values to encrypted fields
		for _, credential := range args.UserEncryptedData {
			for _, field := range credential.EncryptedFields {

				_, err = q.CreateEncryptedData(ctx, CreateEncryptedDataParams{
					FieldName:    field.FieldName,
					FieldValue:   field.FieldValue,
					CredentialID: credential.CredentialID,
					UserID:       args.UserID,
				})
				if err != nil {
					return err
				}
			}
		}

		// Add permissions to access list
		for _, credential := range args.UserEncryptedData {

			accessListParams := AddToAccessListParams{
				CredentialID: credential.CredentialID,
				UserID:       args.UserID,
				AccessType:   credential.AccessType,
				GroupID:      uuid.NullUUID{UUID: args.GroupID, Valid: true},
			}
			_, err = q.AddToAccessList(ctx, accessListParams)
			if err != nil {
				return err
			}

		}

		return nil
	})

	return err
}
