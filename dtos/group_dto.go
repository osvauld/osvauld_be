package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateGroupRequest struct {
	Name string `json:"name"`
}

type GroupDetails struct {
	GroupID   uuid.UUID `json:"groupId"`
	Name      string    `json:"name"`
	CreatedBy uuid.UUID `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
}

type AddMemberToGroupRequest struct {
	GroupID     uuid.UUID                `json:"groupId"`
	MemberID    uuid.UUID                `json:"memberId"`
	MemberRole  string                   `json:"memberRole"`
	Credentials []ShareCredentialPayload `json:"credentials"`
}

type GetUsersOfGroupsRequest struct {
	GroupIDs []uuid.UUID `json:"groupIds"`
}

type RemoveMemberFromGroupRequest struct {
	GroupID  uuid.UUID `json:"groupId"`
	MemberID uuid.UUID `json:"memberId"`
}

type EditGroup struct {
	Name string `json:"name"`
}
