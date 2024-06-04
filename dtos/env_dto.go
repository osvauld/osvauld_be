package dto

import "github.com/google/uuid"

// type EnvField struct {
// 	ID         uuid.UUID `json:"id"`
// 	FieldName  string    `json:"fieldName"`
// 	FieldValue string    `json:"fieldValue"`
// }

type EnvField struct {
	FieldID    uuid.UUID `json:"fieldId"`
	FieldValue string    `json:"fieldValue"`
}

type CredentialsData struct {
	CredentialID uuid.UUID  `json:"credentialId"`
	Fields       []EnvField `json:"fields"`
}

type ShareCredentialsWithEnvironmentRequest struct {
	EnvId       uuid.UUID         `json:"envId"`
	Credentials []CredentialsData `json:"credentials"`
}

type CredentialEnvData struct {
	CredentialID       uuid.UUID
	EnvID              uuid.UUID
	ParentFieldValueID uuid.UUID
	FieldValue         string
	FieldName          string
}

// Todo: later merge this struct and EnvField stuct
type EnvFieldData struct {
	FieldID    uuid.UUID `json:"fieldId"`
	FieldName  string    `json:"fieldName"`
	FieldValue string    `json:"fieldValue"`
}

type CredentialEnvFields struct {
	CredentialID   uuid.UUID      `json:"credentialId"`
	CredentialName string         `json:"credentialName"`
	Fields         []EnvFieldData `json:"fields"`
}
