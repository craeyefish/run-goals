package dto

type CreateGroupMemberRequest struct {
	GroupCode string `json:"group_code"`
	Role      string `json:"role"`
}
