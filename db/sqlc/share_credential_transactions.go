package db

import (
	"context"
	"fmt"
	dto "osvauld/dtos"
	"osvauld/infra/logger"

	"github.com/google/uuid"
)

func (store *SQLStore) ShareCredentialWithUserTransaction(ctx context.Context, args dto.CredentialFieldsForUserDto) error {

	err := store.execTx(ctx, func(q *Queries) error {

		// Create encrypted data records
		for _, field := range args.Fields {
			_, err := q.CreateFieldData(ctx, CreateFieldDataParams{
				FieldName:    field.FieldName,
				FieldValue:   field.FieldValue,
				FieldType:    field.FieldType,
				CredentialID: args.CredentialID,
				UserID:       args.UserID,
			})
			if err != nil {
				return err
			}
		}

		// Add row in access list
		accessListParams := AddToAccessListParams{
			CredentialID: args.CredentialID,
			UserID:       args.UserID,
			AccessType:   args.AccessType,
		}
		_, err := q.AddToAccessList(ctx, accessListParams)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (store *SQLStore) ShareMultipleCredentialsWithMultipleUsersTransaction(ctx context.Context, args []dto.CredentialFieldsForUserDto) error {

	err := store.execTx(ctx, func(q *Queries) error {

		for _, credentialData := range args {

			// Create encrypted data records
			for _, field := range credentialData.Fields {
				_, err := q.CreateFieldData(ctx, CreateFieldDataParams{
					FieldName:    field.FieldName,
					FieldValue:   field.FieldValue,
					CredentialID: credentialData.CredentialID,
					UserID:       credentialData.UserID,
				})
				if err != nil {
					return err
				}
			}

			// Add row in access list
			accessListParams := AddToAccessListParams{
				CredentialID: credentialData.CredentialID,
				UserID:       credentialData.UserID,
				AccessType:   credentialData.AccessType,
			}
			_, err := q.AddToAccessList(ctx, accessListParams)
			if err != nil {
				return err
			}

		}

		return nil
	})

	return err
}

/*
* for each user iterate through all the credentials,
* check if a copy of encrypted fields exists for the user
* if not create a copy of encrypted fields for the user
* else just add the user to the access list
 */
func (store *SQLStore) ShareCredentialWithGroupTransaction(ctx context.Context, groupID uuid.UUID, accessType string, groupPaylod []dto.GroupCredentialPayload) error {
	//TODO: change the payload to be of two types one which just updates the access list and another to add credential and update the access list
	err := store.execTx(ctx, func(q *Queries) error {

		for _, userData := range groupPaylod {
			for _, credential := range userData.Credentials {
				accessLists, err := q.GetCredentialAccessForUser(ctx, GetCredentialAccessForUserParams{
					UserID:       userData.UserID,
					CredentialID: credential.CredentialID,
				})
				if err != nil {
					return err
				}
				if len(accessLists) == 0 {
					for _, field := range credential.EncryptedFields {
						_, err := q.CreateFieldData(ctx, CreateFieldDataParams{
							FieldName:    field.FieldName,
							FieldValue:   field.FieldValue,
							CredentialID: credential.CredentialID,
							UserID:       userData.UserID,
						})
						if err != nil {
							return err
						}
					}
				}
				_, err = q.AddToAccessList(ctx, AddToAccessListParams{
					CredentialID: credential.CredentialID,
					UserID:       userData.UserID,
					GroupID:      uuid.NullUUID{UUID: groupID, Valid: true},
					AccessType:   accessType,
				})
				if err != nil {
					return err
				}

			}
		}

		return nil
	})

	return err
}

func (store *SQLStore) ShareMultipleCredentialsWithMultipleGroupsTransaction(ctx context.Context, args []dto.CredentialEncryptedFieldsForGroupDto) error {

	fmt.Println("started transaction")
	err := store.execTx(ctx, func(q *Queries) error {

		for _, credentialData := range args {
			// Create encrypted data records
			for _, userData := range credentialData.UserEncryptedFields {
				for _, field := range userData.Fields {
					_, err := q.CreateFieldData(ctx, CreateFieldDataParams{
						FieldName:    field.FieldName,
						FieldValue:   field.FieldValue,
						CredentialID: credentialData.CredentialID,
						UserID:       userData.UserID,
					})
					if err != nil {
						return err
					}
				}

				accessListParams := AddToAccessListParams{
					CredentialID: credentialData.CredentialID,
					UserID:       userData.UserID,
					AccessType:   credentialData.AccessType,
					GroupID:      uuid.NullUUID{Valid: true, UUID: credentialData.GroupID},
				}
				_, err := q.AddToAccessList(ctx, accessListParams)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	fmt.Println("ended transaction")

	return err
}

func (store *SQLStore) ShareFolderWithUserTransaction(ctx context.Context, folderId uuid.UUID, credentialPayload dto.CredentialsForUsersPayload) error {

	err := store.execTx(ctx, func(q *Queries) error {

		userId := credentialPayload.UserID
		accessType := credentialPayload.AccessType
		// Create encrypted data records
		for _, credential := range credentialPayload.CredentialData {
			exists, err := q.CheckAccessListEntryExists(ctx, CheckAccessListEntryExistsParams{
				CredentialID: credential.CredentialID,
				UserID:       userId,
			})
			if err != nil {
				return err
			}
			if !exists {
				for _, field := range credential.EncryptedFields {
					_, err = q.CreateFieldData(ctx, CreateFieldDataParams{
						FieldName:    field.FieldName,
						FieldValue:   field.FieldValue,
						CredentialID: credential.CredentialID,
						UserID:       userId,
					})
					if err != nil {
						logger.Debugf("\nerror: %v ", err)
						return err
					}
				}
			}
			_, err = q.AddToAccessList(ctx, AddToAccessListParams{
				CredentialID: credential.CredentialID,
				UserID:       userId,
				AccessType:   accessType,
			})
			if err != nil {
				return err
			}
		}
		q.AddFolderAccess(ctx, AddFolderAccessParams{
			FolderID:   folderId,
			UserID:     userId,
			AccessType: accessType,
		})

		return nil
	})

	return err
}

func (store *SQLStore) ShareFolderWithGroupTransaction(ctx context.Context, folderId uuid.UUID, credentialPayloads dto.CredentialsForGroupsPayload) error {

	err := store.execTx(ctx, func(q *Queries) error {

		groupId := credentialPayloads.GroupID
		accessType := credentialPayloads.AccessType
		for _, userData := range credentialPayloads.EncryptedUserData {
			for _, credentialData := range userData.Credentials {
				exists, err := q.CheckAccessListEntryExists(ctx, CheckAccessListEntryExistsParams{
					CredentialID: credentialData.CredentialID,
					UserID:       userData.UserID,
				})
				if err != nil {
					return err
				}
				if !exists {
					for _, field := range credentialData.EncryptedFields {
						_, err = q.CreateFieldData(ctx, CreateFieldDataParams{
							FieldName:    field.FieldName,
							FieldValue:   field.FieldValue,
							CredentialID: credentialData.CredentialID,
							UserID:       userData.UserID,
						})
						if err != nil {
							return err
						}
					}
				}
				_, err = q.AddToAccessList(ctx, AddToAccessListParams{
					CredentialID: credentialData.CredentialID,
					UserID:       userData.UserID,
					AccessType:   accessType,
					GroupID:      uuid.NullUUID{Valid: true, UUID: groupId},
				})
				if err != nil {
					return err
				}
			}

			err := q.AddFolderAccessWithGroup(ctx, AddFolderAccessWithGroupParams{
				FolderID:   folderId,
				UserID:     userData.UserID,
				AccessType: accessType,
				GroupID:    uuid.NullUUID{Valid: true, UUID: groupId},
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}
