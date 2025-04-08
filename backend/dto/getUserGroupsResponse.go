package dto

import (
	"run-goals/models"
)

type GetUserGroupsResponse struct {
	Groups    []models.Group `json:"groups"`
	Name      string         `json:"name"`
	CreatedBy int64          `json:"created_by"`
}
