package dto

type CreateGroupMemberRequest struct {
	GroupID int64  `json:"group_id"`
	UserID  int64  `json:"user_id"`
	Role    string `json:"role"`
}
