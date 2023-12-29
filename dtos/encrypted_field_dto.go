package dto

import "github.com/google/uuid"

type CredentialEncryptedFieldsDto struct {
	CredentialID    uuid.UUID `json:"credentialId"`
	EncryptedFields []Field   `json:"encryptedFields"`
	AccessType      string    `json:"accessType"`
}
