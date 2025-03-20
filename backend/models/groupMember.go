package models

import "time"

type GroupMember struct {
	ID       int64     `json:"id"`
	GroupId  int64     `json:"group_id"`
	UserId   int64     `json:"user_id"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}
