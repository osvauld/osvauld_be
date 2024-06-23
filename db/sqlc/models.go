// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Credential struct {
	ID             uuid.UUID      `json:"id"`
	Name           string         `json:"name"`
	Description    sql.NullString `json:"description"`
	CredentialType string         `json:"credentialType"`
	FolderID       uuid.UUID      `json:"folderId"`
	Domain         sql.NullString `json:"domain"`
	CreatedBy      uuid.NullUUID  `json:"createdBy"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedBy      uuid.NullUUID  `json:"updatedBy"`
	UpdatedAt      time.Time      `json:"updatedAt"`
}

type CredentialAccess struct {
	ID           uuid.UUID     `json:"id"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
	CredentialID uuid.UUID     `json:"credentialId"`
	UserID       uuid.UUID     `json:"userId"`
	AccessType   string        `json:"accessType"`
	GroupID      uuid.NullUUID `json:"groupId"`
	FolderID     uuid.NullUUID `json:"folderId"`
}

type Environment struct {
	ID        uuid.UUID `json:"id"`
	CliUser   uuid.UUID `json:"cliUser"`
	Name      string    `json:"name"`
	Createdat time.Time `json:"createdat"`
	Updatedat time.Time `json:"updatedat"`
	CreatedBy uuid.UUID `json:"createdBy"`
}

type EnvironmentField struct {
	ID                 uuid.UUID `json:"id"`
	FieldName          string    `json:"fieldName"`
	FieldValue         string    `json:"fieldValue"`
	ParentFieldValueID uuid.UUID `json:"parentFieldValueId"`
	EnvID              uuid.UUID `json:"envId"`
	CredentialID       uuid.UUID `json:"credentialId"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

type Field struct {
	ID           uuid.UUID     `json:"id"`
	FieldName    string        `json:"fieldName"`
	FieldValue   string        `json:"fieldValue"`
	FieldType    string        `json:"fieldType"`
	CredentialID uuid.UUID     `json:"credentialId"`
	UserID       uuid.UUID     `json:"userId"`
	CreatedAt    time.Time     `json:"createdAt"`
	CreatedBy    uuid.NullUUID `json:"createdBy"`
	UpdatedAt    time.Time     `json:"updatedAt"`
	UpdatedBy    uuid.NullUUID `json:"updatedBy"`
}

type FieldDatum struct {
	ID           uuid.UUID     `json:"id"`
	FieldName    string        `json:"fieldName"`
	FieldType    string        `json:"fieldType"`
	CredentialID uuid.UUID     `json:"credentialId"`
	CreatedAt    time.Time     `json:"createdAt"`
	CreatedBy    uuid.NullUUID `json:"createdBy"`
	UpdatedAt    time.Time     `json:"updatedAt"`
	UpdatedBy    uuid.NullUUID `json:"updatedBy"`
}

type FieldValue struct {
	ID         uuid.UUID `json:"id"`
	FieldID    uuid.UUID `json:"fieldId"`
	FieldValue string    `json:"fieldValue"`
	UserID     uuid.UUID `json:"userId"`
}

type Folder struct {
	ID          uuid.UUID      `json:"id"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
	CreatedBy   uuid.NullUUID  `json:"createdBy"`
	Type        string         `json:"type"`
}

type FolderAccess struct {
	ID         uuid.UUID     `json:"id"`
	CreatedAt  time.Time     `json:"createdAt"`
	UpdatedAt  time.Time     `json:"updatedAt"`
	FolderID   uuid.UUID     `json:"folderId"`
	UserID     uuid.UUID     `json:"userId"`
	AccessType string        `json:"accessType"`
	GroupID    uuid.NullUUID `json:"groupId"`
}

type GroupList struct {
	ID         uuid.UUID `json:"id"`
	GroupingID uuid.UUID `json:"groupingId"`
	UserID     uuid.UUID `json:"userId"`
	AccessType string    `json:"accessType"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Grouping struct {
	ID        uuid.UUID     `json:"id"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
	Name      string        `json:"name"`
	CreatedBy uuid.NullUUID `json:"createdBy"`
}

type SessionTable struct {
	ID        uuid.UUID      `json:"id"`
	UserID    uuid.UUID      `json:"userId"`
	PublicKey string         `json:"publicKey"`
	Challenge string         `json:"challenge"`
	DeviceID  sql.NullString `json:"deviceId"`
	SessionID sql.NullString `json:"sessionId"`
	CreatedAt sql.NullTime   `json:"createdAt"`
	UpdatedAt sql.NullTime   `json:"updatedAt"`
}

type User struct {
	ID                    uuid.UUID      `json:"id"`
	CreatedAt             time.Time      `json:"createdAt"`
	UpdatedAt             time.Time      `json:"updatedAt"`
	Username              string         `json:"username"`
	Name                  string         `json:"name"`
	EncryptionKey         sql.NullString `json:"encryptionKey"`
	DeviceKey             sql.NullString `json:"deviceKey"`
	TempPassword          string         `json:"tempPassword"`
	RegistrationChallenge sql.NullString `json:"registrationChallenge"`
	SignedUp              bool           `json:"signedUp"`
	Type                  string         `json:"type"`
	Status                string         `json:"status"`
	CreatedBy             uuid.NullUUID  `json:"createdBy"`
}
