package db

import (
	"context"
	dto "osvauld/dtos"
)

func (store *SQLStore) AddCredentialFieldToEnvTxn(ctx context.Context, args []dto.CredentialEnvData) error {

	err := store.execTx(ctx, func(q *Queries) error {
		for _, arg := range args {
			_, err := q.CreateEnvFields(ctx, CreateEnvFieldsParams{
				EnvID:              arg.EnvID,
				ParentFieldValueID: arg.ParentFieldValueID,
				FieldValue:         arg.FieldValue,
				FieldName:          arg.FieldName,
				CredentialID:       arg.CredentialID,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err

}
