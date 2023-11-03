package repository

import (
	db "osvauld/db/sqlc"
	dto "osvauld/dtos"
	"osvauld/infra/database"
	"osvauld/infra/logger"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddGroup(ctx *gin.Context, group dto.CreateGroup, userID uuid.UUID) error {
	arg := db.CreateGroupParams{
		Name:      group.Name,
		CreatedBy: uuid.NullUUID{UUID: userID, Valid: true},
		Members:   []uuid.UUID{userID},
	}
	q := db.New(database.DB)
	_, err := q.CreateGroup(ctx, arg)
	return err
}

func AddMembersToGroup(ctx *gin.Context, payload dto.AddMembers, userID uuid.UUID) error {
	uuidStrings := make([]string, len(payload.Members))
	for i, u := range payload.Members {
		uuidStrings[i] = u.String()
	}
	pgArray := "{" + strings.Join(uuidStrings, ",") + "}"
	arg := db.AddMemberToGroupParams{
		CreatedBy:   uuid.NullUUID{UUID: userID, Valid: true},
		ArrayAppend: pgArray,
		ID:          payload.GroupID,
	}
	q := db.New(database.DB)
	err := q.AddMemberToGroup(ctx, arg)
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	return nil

}
