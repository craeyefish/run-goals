package responses

import "run-goals/models"

type ListPeakSummariesResponse struct {
	PeakSummaries []models.PeakSummary `json:"peak_summaries"`
}
