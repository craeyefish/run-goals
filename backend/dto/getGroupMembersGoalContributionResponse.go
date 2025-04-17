package dto

import (
	"run-goals/models"
)

type GetGroupMembersGoalContributionResponse struct {
	Members []models.GroupMemberGoalContribution `json:"members"`
}
