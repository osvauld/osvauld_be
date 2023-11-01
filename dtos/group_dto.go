package dto

import "github.com/google/uuid"

type CreateGroup struct {
	Name string `json:"name"`
}

type AddMembers struct {
	Members []uuid.UUID `json:"members"`
	GroupID uuid.UUID   `json:"groupId"`
}
