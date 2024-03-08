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
	UserIDs []uuid.UUID `json:"userIds"`
}

type RemoveFolderAccessForUsers struct {
	UserIDs []uuid.UUID `json:"userIds"`
}

type RemoveCredentialAccessForGroups struct {
	GroupIDs []uuid.UUID `json:"groupIds"`
}

type RemoveFolderAccessForGroups struct {
	GroupIDs []uuid.UUID `json:"groupIds"`
}
