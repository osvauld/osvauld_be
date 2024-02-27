// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

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
	CreatedBy      uuid.UUID      `json:"createdBy"`
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

type Field struct {
	ID           uuid.UUID     `json:"id"`
	FieldName    string        `json:"fieldName"`
	FieldValue   string        `json:"fieldValue"`
	FieldType    string        `json:"fieldType"`
	CredentialID uuid.UUID     `json:"credentialId"`
	UserID       uuid.UUID     `json:"userId"`
	CreatedAt    time.Time     `json:"createdAt"`
	CreatedBy    uuid.UUID     `json:"createdBy"`
	UpdatedAt    time.Time     `json:"updatedAt"`
	UpdatedBy    uuid.NullUUID `json:"updatedBy"`
}

type FieldArchive struct {
	ID         uuid.UUID `json:"id"`
	FieldID    uuid.UUID `json:"fieldId"`
	FieldName  string    `json:"fieldName"`
	FieldValue string    `json:"fieldValue"`
	FieldType  string    `json:"fieldType"`
	CreateAt   time.Time `json:"createAt"`
	CreatedBy  uuid.UUID `json:"createdBy"`
	UpdatedAt  time.Time `json:"updatedAt"`
	UpdatedBy  uuid.UUID `json:"updatedBy"`
	Version    int32     `json:"version"`
}

type Folder struct {
	ID          uuid.UUID      `json:"id"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
	CreatedBy   uuid.UUID      `json:"createdBy"`
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
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string    `json:"name"`
	CreatedBy uuid.UUID `json:"createdBy"`
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
	Status                string         `json:"status"`
}
