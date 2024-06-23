package dto

import "github.com/google/uuid"

type UserFields struct {
	UserID uuid.UUID `json:"userId"`
	Fields []Field   `json:"fields"`
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

type Fields struct {
	FieldID     uuid.UUID    `json:"fieldId"`
	FieldName   string       `json:"fieldName"`
	FieldType   string       `json:"fieldType"`
	FieldValues []FieldValue `json:"fieldValues"`
}


type FieldValue struct {
	UserID     uuid.UUID `json:"userId"`
	FieldValue string    `json:"fieldValue"`
}
