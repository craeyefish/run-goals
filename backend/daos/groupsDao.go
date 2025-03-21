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
	DeleteGroupMember(userID int64) error
	CreateGroupGoal(goal models.GroupGoal) error
	UpdateGroupGoal(goal models.GroupGoal) error
	DeleteGroupGoal(goalID int64) error

	GetUserGroups(userID int64) ([]models.Group, error)
	GetGroupMembers(groupID int64) ([]models.GroupMember, error)
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
	_, err := dao.db.Exec(sql, member.GroupID, member.UserID, member.Role, member.JoinedAt)
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

func (dao *GroupsDao) DeleteGroupMember(userID int64) error {
	sql := `
		DELETE FROM
			group_members
		WHERE id = $1;
	`
	_, err := dao.db.Exec(sql, userID)
	if err != nil {
		dao.l.Printf("Error deleting group member: %v", err)
		return err
	}
	return nil
}

func (dao *GroupsDao) CreateGroupGoal(goal models.GroupGoal) error {
	sql := `
		INSERT INTO group_goals (
			group_id,
			name,
			target_value,
			start_date,
			end_date,
			created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6
		);
	`
	_, err := dao.db.Exec(sql, goal.GroupID, goal.Name, goal.TargetValue, goal.StartDate, goal.EndDate, goal.CreatedAt)
	if err != nil {
		dao.l.Printf("Error adding group goal: %v", err)
		return err
	}
	return nil
}

func (dao *GroupsDao) UpdateGroupGoal(goal models.GroupGoal) error {
	sql := `
		UPDATE
			group_goals
		SET
			name = $2,
			target_vale = $3,
			start_date = $4,
			end_date = $5,
			created_at = $6
		WHERE id = $1;
	`
	_, err := dao.db.Exec(sql, goal.ID, goal.Name, goal.TargetValue, goal.StartDate, goal.EndDate, goal.CreatedAt)
	if err != nil {
		dao.l.Printf("Error updating group goal: %v", err)
		return err
	}
	return nil
}

func (dao *GroupsDao) DeleteGroupGoal(goalID int64) error {
	sql := `
		DELETE FROM
			group_goals
		WHERE id = $1;
	`
	_, err := dao.db.Exec(sql, goalID)
	if err != nil {
		dao.l.Printf("Error deleting group goal: %v", err)
		return err
	}
	return nil
}

func (dao *GroupsDao) GetUserGroups(userID int64) ([]models.Group, error) {
	groups := []models.Group{}
	sql := `
		SELECT
			id,
			name,
			created_by,
			created_at
		FROM groups
		WHERE
			id = $1;
	`
	rows, err := dao.db.Query(sql, userID)
	if err != nil {
		dao.l.Printf("Error getting user groups: %v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		group := models.Group{}
		err = rows.Scan(
			&group.ID,
			&group.Name,
			&group.CreatedBy,
			&group.CreatedAt,
		)
		if err != nil {
			dao.l.Printf("Error parsing query result: %f", err)
			return nil, err
		}
		groups = append(groups, group)
	}
	err = rows.Err()
	if err != nil {
		dao.l.Printf("Error during iteration: %v", err)
		return nil, err
	}

	return groups, nil
}

func (dao *GroupsDao) GetGroupMembers(groupID int64) ([]models.GroupMember, error) {
	groupMembers := []models.GroupMember{}
	sql := `
		SELECT
			user_id
		FROM group_members
		WHERE group_id = $1;
	`
	rows, err := dao.db.Query(sql, groupID)
	if err != nil {
		dao.l.Printf("Error getting group members: %v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		groupMember := models.GroupMember{}
		err = rows.Scan(
			&groupMember.ID,
			&groupMember.GroupID,
			&groupMember.UserID,
			&groupMember.Role,
			&groupMember.JoinedAt,
		)
		if err != nil {
			dao.l.Printf("Error parsing query result: %f", err)
			return nil, err
		}
		groupMembers = append(groupMembers, groupMember)
	}
	err = rows.Err()
	if err != nil {
		dao.l.Printf("Error during iteration: %v", err)
		return nil, err
	}

	return groupMembers, nil
}
