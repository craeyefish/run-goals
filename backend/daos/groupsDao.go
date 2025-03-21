package daos

import (
	"database/sql"
	"log"
	"run-goals/models"
)

type GroupsDaoInterface interface {
	CreateGroup(group models.Group) error
	UpdateGroup(group models.Group) error
	DeleteGroup(groupID int64) error
	CreateGroupMember(member models.GroupMember) error
	UpdateGroupMember(member models.GroupMember) error
	DeleteGroupMember(userID int64) error        // todo
	CreateGroupGoal(goal models.GroupGoal) error // todo
	UpdateGroupGoal(goal models.GroupGoal) error // todo
	DeleteGroupGoal(goalID int64) error          // todo

	GetUserGroups(userID int64) ([]models.Group, error)          // todo
	GetGroupMembers(groupID int64) ([]models.GroupMember, error) // todo
}

type GroupsDao struct {
	l  *log.Logger
	db *sql.DB
}

func NewGroupsDao(logger *log.Logger, db *sql.DB) *GroupsDao {
	return &GroupsDao{
		l:  logger,
		db: db,
	}
}

func (dao *GroupsDao) CreateGroup(group models.Group) error {
	sql := `
		INSERT INTO groups (
			name,
			created_by,
			created_at
		) VALUES (
			$1, $2, $3
		);
	`
	_, err := dao.db.Exec(sql, group.Name, group.CreatedBy, group.CreatedAt)
	if err != nil {
		dao.l.Printf("Error creating group: %v", err)
		return err
	}
	return nil
}

func (dao *GroupsDao) UpdateGroup(group models.Group) error {
	sql := `
		UPDATE
			groups
		SET
			name = $2
		WHERE id = $1;
	`
	_, err := dao.db.Exec(sql, group.ID, group.Name)
	if err != nil {
		dao.l.Printf("Error updating group: %v", err)
		return err
	}
	return nil
}

func (dao *GroupsDao) DeleteGroup(groupID int64) error {
	sql := `
		DELETE FROM
			groups
		WHERE id = $1;
	`
	_, err := dao.db.Exec(sql, groupID)
	if err != nil {
		dao.l.Printf("Error deleting group: %v", err)
		return err
	}
	return nil
}

func (dao *GroupsDao) CreateGroupMember(member models.GroupMember) error {
	sql := `
		INSERT INTO group_members (
			group_id,
			user_id,
			role,
			joined_at
		) VALUES (
			$1, $2, $3, $4
		);
	`
	_, err := dao.db.Exec(sql, member.GroupId, member.UserId, member.Role, member.JoinedAt)
	if err != nil {
		dao.l.Printf("Error adding group member: %v", err)
		return err
	}
	return nil
}

func (dao *GroupsDao) UpdateGroupMember(member models.GroupMember) error {
	sql := `
		UPDATE
			group_members
		SET
			role = $2
		WHERE id = $1;
	`
	_, err := dao.db.Exec(sql, member.ID, member.Role)
	if err != nil {
		dao.l.Printf("Error updating group member: %v", err)
		return err
	}
	return nil
}
