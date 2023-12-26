package dto

import "github.com/google/uuid"

type CredentialEncryptedFielsdDto struct {
	CredentialID    uuid.UUID `json:"credentialId"`
	EncryptedFields []Field   `json:"encryptedFields"`
}
