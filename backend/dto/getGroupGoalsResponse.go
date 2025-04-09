package dto

import (
	"run-goals/models"
)

type GetGroupGoalsResponse struct {
	Goals []models.GroupGoal `json:"goals"`
}
