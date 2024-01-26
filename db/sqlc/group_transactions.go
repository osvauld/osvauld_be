package db

import (
	"context"
	dto "osvauld/dtos"

	"github.com/google/uuid"
)

func (store *SQLStore) CreateGroupAndAddManager(ctx context.Context, groupData dto.GroupDetails) (dto.GroupDetails, error) {

	var createGroupResult CreateGroupRow
	err := store.execTx(ctx, func(q *Queries) error {
		createGroupResult, err := q.CreateGroup(ctx, CreateGroupParams{
			Name:      groupData.Name,
			CreatedBy: groupData.CreatedBy,
		})
		if err != nil {
			return err
		}

		err = q.AddGroupMember(ctx, AddGroupMemberParams{
			GroupingID: createGroupResult.ID,
			UserID:     groupData.CreatedBy,
			AccessType: "manager",
		})
		if err != nil {
			return err
		}

		return nil
	})

	groupData.GroupID = createGroupResult.ID
	groupData.CreatedAt = createGroupResult.CreatedAt

	return groupData, err
}

type AddMemberToGroupTransactionParams struct {
	GroupID           uuid.UUID                        `json:"groupId"`
	UserID            uuid.UUID                        `json:"userId"`
	MemberRole        string                           `json:"memberRole"`
	UserEncryptedData []dto.CredentialFieldsForUserDto `json:"encryptedFields"`
}

func (store *SQLStore) AddMemberToGroupTransaction(ctx context.Context, args AddMemberToGroupTransactionParams) error {
	err := store.execTx(ctx, func(q *Queries) error {

		// Add record to grouping table
		err := q.AddGroupMember(ctx, AddGroupMemberParams{
			GroupingID: args.GroupID,
			UserID:     args.UserID,
			AccessType: args.MemberRole,
		})
		if err != nil {
			return err
		}

		// Add Values to encrypted fields
		for _, credential := range args.UserEncryptedData {
			for _, field := range credential.Fields {

				_, err = q.CreateFieldData(ctx, CreateFieldDataParams{
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
