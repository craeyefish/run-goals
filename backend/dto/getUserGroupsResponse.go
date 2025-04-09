package dto

import (
	"run-goals/models"
)

type GetUserGroupsResponse struct {
	Groups []models.Group `json:"groups"`
}
