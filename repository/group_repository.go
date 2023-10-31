package repository

import (
	"osvauld/infra/database"
	"osvauld/models"

	"github.com/google/uuid"
)

func AddGroup(group models.Group) error {
	db := database.DB
	return db.Create(&group).Error
}

func AddMembersToGroup(groupID uuid.UUID, members []uuid.UUID) error {
	db := database.DB

	// Convert the slice of UUIDs to a string format suitable for PostgreSQL
	idsStr := "{"
	for i, member := range members {
		idsStr += member.String()
		if i < len(members)-1 {
			idsStr += ","
		}
	}
	idsStr += "}"

	// Raw SQL query to append userIDs to the members array
	query := `UPDATE groups SET members = array_cat(members, ?::uuid[]) WHERE id = ?`

	return db.Exec(query, idsStr, groupID).Error
}
