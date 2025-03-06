package models

type GoalProgress struct {
	Goal            float64            `json:"goal"` // total km goal
	CurrentProgress float64            `json:"currentProgress"`
	Contributions   []UserContribution `json:"contributions"`
}
