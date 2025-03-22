package dto

type CreateGroupRequest struct {
	Name      string `json:"name"`
	CreatedBy int64  `json:"created_by"`
}
