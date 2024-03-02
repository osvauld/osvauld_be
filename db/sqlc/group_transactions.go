package db

import (
	"context"
	dto "osvauld/dtos"
)

func (store *SQLStore) CreateGroupAndAddManager(ctx context.Context, groupData dto.GroupDetails) (dto.GroupDetails, error) {

	var createGroupResult CreateGroupRow
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		createGroupResult, err = q.CreateGroup(ctx, CreateGroupParams{
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

type AddMembersToGroupTransactionParams struct {
	FieldArgs            []AddFieldParams            `json:"fieldArgs"`
	CredentialAccessArgs []AddCredentialAccessParams `json:"credentialAccessArgs"`
	FolderAccessArgs     []AddFolderAccessParams     `json:"folderAccessArgs"`
	GroupMembershipArgs  []AddGroupMemberParams      `json:"groupMembershipArgs"`
}

func (store *SQLStore) AddMembersToGroupTransaction(ctx context.Context, args AddMembersToGroupTransactionParams) error {

	err := store.execTx(ctx, func(q *Queries) error {

		for _, fieldRecord := range args.FieldArgs {

			_, err := q.AddField(ctx, fieldRecord)
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
