package dto

import "github.com/google/uuid"

type AccessListResult struct {
	ID           uuid.UUID     `json:"id"`
	UserID       uuid.UUID     `json:"user_id"`
	CredentialID uuid.UUID     `json:"credential_id"`
	GroupID      uuid.NullUUID `json:"group_id"`
	AccessType   string        `json:"access_type"`
}

type RemoveCredentialAccessForUsers struct {
	CredentialID uuid.UUID   `json:"credential_id"`
	UserIDs      []uuid.UUID `json:"user_ids"`
}
