package dto

import (
	"run-goals/models"
)

type GetGroupMembersResponse struct {
	Members []models.GroupMember `json:"members"`
}
