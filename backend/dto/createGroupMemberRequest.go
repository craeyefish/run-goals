package dto

type CreateGroupMemberRequest struct {
	GroupCode string `json:"group_code"`
	UserID    int64  `json:"user_id"`
	Role      string `json:"role"`
}
