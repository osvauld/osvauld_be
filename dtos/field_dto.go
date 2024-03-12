package dto

import "github.com/google/uuid"

type UserFields struct {
	UserID uuid.UUID `json:"userId"`
	Fields []Field   `json:"fields"`
}

type UserFieldsWithAccessType struct {
	UserID     uuid.UUID `json:"userId"`
	Fields     []Field   `json:"fields"`
	AccessType string    `json:"accessType"`
}

type Field struct {
	ID         uuid.UUID `json:"fieldId"`
	FieldName  string    `json:"fieldName"`
	FieldValue string    `json:"fieldValue"`
	FieldType  string    `json:"fieldType"`
}

type ShareField struct {
	ID         uuid.UUID `json:"fieldId"`
	FieldValue string    `json:"fieldValue"`
}

type CredentialFields struct {
	CredentialID uuid.UUID `json:"credentialId"`
	Fields       []Field   `json:"fields"`
}