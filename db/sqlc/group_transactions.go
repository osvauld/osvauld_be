package db

import (
	"context"
	dto "osvauld/dtos"

	"github.com/google/uuid"
)

func (store *SQLStore) CreateGroupAndAddManager(ctx context.Context, groupData dto.GroupDetails) (dto.GroupDetails, error) {

	var createGroupResult CreateGroupRow
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		createGroupResult, err = q.CreateGroup(ctx, CreateGroupParams{
			Name:      groupData.Name,
			CreatedBy: uuid.NullUUID{UUID: groupData.CreatedBy, Valid: true},
		})
		if err != nil {
			return err
		}

		err = q.AddGroupMember(ctx, AddGroupMemberParams{
			GroupingID: createGroupResult.ID,
			UserID:     groupData.CreatedBy,
			AccessType: "admin",
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

type AddMembersToGroupTransactionParams struct {
	FieldArgs            []AddFieldValueParams       `json:"fieldArgs"`
	CredentialAccessArgs []AddCredentialAccessParams `json:"credentialAccessArgs"`
	FolderAccessArgs     []AddFolderAccessParams     `json:"folderAccessArgs"`
	GroupMembershipArgs  []AddGroupMemberParams      `json:"groupMembershipArgs"`
}

func (store *SQLStore) AddMembersToGroupTransaction(ctx context.Context, args AddMembersToGroupTransactionParams) error {

	err := store.execTx(ctx, func(q *Queries) error {

		for _, fieldRecord := range args.FieldArgs {

			_, err := q.AddFieldValue(ctx, fieldRecord)
			if err != nil {
				return err
			}
		}

		for _, credentialAccessRecord := range args.CredentialAccessArgs {

			_, err := q.AddCredentialAccess(ctx, credentialAccessRecord)
			if err != nil {
				return err
			}

		}

		for _, folderAccessRecord := range args.FolderAccessArgs {

			err := q.AddFolderAccess(ctx, folderAccessRecord)
			if err != nil {
				return err
			}

		}

		for _, groupMembershipRecord := range args.GroupMembershipArgs {

			err := q.AddGroupMember(ctx, groupMembershipRecord)
			if err != nil {
				return err
			}

		}

		return nil
	})

	return err
}

type RemoveMemberFromGroupTransactionParams struct {
	GroupID  uuid.UUID
	MemberID uuid.UUID
}

func (store *SQLStore) RemoveMemberFromGroupTransaction(ctx context.Context, args RemoveMemberFromGroupTransactionParams) error {

	err := store.execTx(ctx, func(q *Queries) error {

		err := q.RemoveUserFromGroupList(ctx, RemoveUserFromGroupListParams{
			GroupingID: args.GroupID,
			UserID:     args.MemberID,
		})
		if err != nil {
			return err
		}

		err = q.RemoveCredentialAccessForGroupMember(ctx, RemoveCredentialAccessForGroupMemberParams{
			UserID:  args.MemberID,
			GroupID: uuid.NullUUID{UUID: args.GroupID, Valid: true},
		})
		if err != nil {
			return err
		}

		err = q.RemoveFolderAccessForGroupMember(ctx, RemoveFolderAccessForGroupMemberParams{
			UserID:  args.MemberID,
			GroupID: uuid.NullUUID{UUID: args.GroupID, Valid: true},
		})
		if err != nil {
			return err
		}

		return nil

	})

	return err
}
