package responses

import "run-goals/models"

type ListPeakResponse struct {
	models.Peak
	IsSummited bool `json:"is_summited"`
}
