package daos

import (
	"database/sql"
	"log"
	"run-goals/models"
	"time"

	"github.com/lib/pq"
)

type GroupsDaoInterface interface {
	CreateGroup(request models.Group) (int64, error)
	UpdateGroup(group models.Group) error
	DeleteGroup(groupID int64) error
	CheckGroupCodeExists(code string) (int64, error)

	CreateGroupMember(member models.GroupMember) error
	UpdateGroupMember(member models.GroupMember) error
	DeleteGroupMember(userID int64) error
	GetGroupMembersGoalContribution(groupID int64, startDate time.Time, endDate time.Time) ([]models.GroupMemberGoalContribution, error)

	CreateGroupGoal(goal models.GroupGoal) (int64, error)
	UpdateGroupGoal(goal models.GroupGoal) error
	DeleteGroupGoal(goalID int64) error

	GetUserGroups(userID int64) ([]models.Group, error)
	GetGroupMembers(groupID int64) ([]models.GroupMember, error)
	GetGroupGoals(groupID int64) ([]models.GroupGoal, error)
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

func (dao *GroupsDao) CreateGroup(group models.Group) (*int64, error) {
	var id int64
	sql := `
		INSERT INTO groups (
			name,
			code,
			created_by,
			created_at
		) VALUES (
			$1, $2, $3, $4
		)
		RETURNING id;
	`
	err := dao.db.QueryRow(sql, group.Name, group.Code, group.CreatedBy, group.CreatedAt).Scan(&id)
	if err != nil {
		dao.l.Printf("Error creating group: %v", err)
		return nil, err
	}
	return &id, nil
}

func (dao *GroupsDao) UpdateGroup(groupID int64, name string) error {
	sql := `
		UPDATE
			groups
		SET
			name = $2
		WHERE id = $1;
	`
	_, err := dao.db.Exec(sql, groupID, name)
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

func (dao *GroupsDao) GetGroupIDFromCode(code string) (*int64, error) {
	var id int64
	sql := `
		SELECT
			id
		FROM groups
		WHERE code = $1;
	`
	row := dao.db.QueryRow(sql, code)
	err := row.Scan(&id)
	if err != nil {
		dao.l.Printf("Error getting group id: %v", err)
		return nil, err
	}

	return &id, nil
}

func (dao *GroupsDao) CheckGroupCodeExists(code string) (*int64, error) {
	var count int64
	sql := `
		SELECT
			COUNT(code)
	 	FROM
			groups
		WHERE code = $1;
	`
	row := dao.db.QueryRow(sql, code)
	err := row.Scan(&count)
	if err != nil {
		dao.l.Println("Error checking group code", err)
		return nil, err
	}

	return &count, nil
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

func (dao *GroupsDao) UpdateGroupMember(groupID int64, userID int64, role string) error {
	sql := `
		UPDATE
			group_members
		SET
			role = $1
		WHERE
			group_id = $2
			AND user_id = $3;
	`
	_, err := dao.db.Exec(sql, role, groupID, userID)
	if err != nil {
		dao.l.Printf("Error updating group member: %v", err)
		return err
	}
	return nil
}

func (dao *GroupsDao) DeleteGroupMember(userID int64, groupID int64) error {
	sql := `
		DELETE FROM
			group_members
		WHERE
			user_id = $1
			AND group_id = $2;
	`
	_, err := dao.db.Exec(sql, userID, groupID)
	if err != nil {
		dao.l.Printf("Error deleting group member: %v", err)
		return err
	}
	return nil
}

func (dao *GroupsDao) CreateGroupGoal(goal models.GroupGoal) (*int64, error) {
	var id int64
	sql := `
        INSERT INTO group_goals (
            group_id,
            name,
            description,
            goal_type,
            target_value,
            target_summits,
            start_date,
            end_date,
            created_at
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9
        )
        RETURNING id;
    `
	err := dao.db.QueryRow(
		sql,
		goal.GroupID,
		goal.Name,
		goal.Description,
		goal.GoalType,
		goal.TargetValue,
		pq.Array(goal.TargetSummits),
		goal.StartDate,
		goal.EndDate,
		goal.CreatedAt,
	).Scan(&id)
	if err != nil {
		dao.l.Printf("Error adding group goal: %v", err)
		return nil, err
	}
	return &id, nil
}

func (dao *GroupsDao) UpdateGroupGoal(goal models.GroupGoal) error {
	sql := `
        UPDATE
            group_goals
        SET
            name = $2,
            description = $3,
            goal_type = $4,
            target_value = $5,
            target_summits = $6,
            start_date = $7,
            end_date = $8
        WHERE id = $1;
    `
	_, err := dao.db.Exec(
		sql,
		goal.ID,
		goal.Name,
		goal.Description,
		goal.GoalType,
		goal.TargetValue,
		pq.Array(goal.TargetSummits),
		goal.StartDate,
		goal.EndDate,
	)
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
			code,
			created_by,
			created_at
		FROM groups
		WHERE
			id in (
				SELECT
					group_id
				FROM group_members
				WHERE user_id = $1
			)
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
			&group.Code,
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
			id,
			group_id,
			user_id,
			role,
			joined_at
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

func (dao *GroupsDao) GetGroupGoals(groupID int64) ([]models.GroupGoal, error) {
	groupGoals := []models.GroupGoal{}
	sql := `
        SELECT
            id,
            group_id,
            name,
            description,
            goal_type,
            target_value,
            target_summits,
            start_date,
            end_date,
            created_at
        FROM group_goals
        WHERE group_id = $1;
    `
	rows, err := dao.db.Query(sql, groupID)
	if err != nil {
		dao.l.Printf("Error getting group goals: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		groupGoal := models.GroupGoal{}
		var targetSummits pq.Int64Array

		err = rows.Scan(
			&groupGoal.ID,
			&groupGoal.GroupID,
			&groupGoal.Name,
			&groupGoal.Description,
			&groupGoal.GoalType,
			&groupGoal.TargetValue,
			&targetSummits,
			&groupGoal.StartDate,
			&groupGoal.EndDate,
			&groupGoal.CreatedAt,
		)
		if err != nil {
			dao.l.Printf("Error parsing query result: %v", err)
			return nil, err
		}

		// Convert pq.Int64Array to []int64
		groupGoal.TargetSummits = []int64(targetSummits)

		groupGoals = append(groupGoals, groupGoal)
	}

	err = rows.Err()
	if err != nil {
		dao.l.Printf("Error during iteration: %v", err)
		return nil, err
	}

	return groupGoals, nil
}

func (dao *GroupsDao) GetGroupMembersGoalContribution(groupID int64, startDate time.Time, endDate time.Time) ([]models.GroupMemberGoalContribution, error) {
	groupMembersGoalContribution := []models.GroupMemberGoalContribution{}
	sql := `
		WITH members_tbl AS (
			SELECT
				id as group_member_id,
				group_id,
				user_id,
				role,
				joined_at
			FROM group_members
			WHERE group_id = $1
		),

		member_activity_tbl AS (
			SELECT
				user_id,
				count(id) as total_activities,
				sum(distance) as total_distance
			FROM activity
			WHERE user_id in (SELECT user_id FROM members_tbl)
				AND start_date >= $2
				AND start_date <= $3
			GROUP BY user_id
		),

		member_peaks AS (
			SELECT
				user_id,
				count(distinct peak_id) as total_unique_summits,
				count(peak_id) as total_summits
			FROM user_peaks
			WHERE user_id in (SELECT user_id FROM members_tbl)
				AND summited_at >= $2
				AND summited_at <= $3
			GROUP BY user_id
		)

		SELECT
			members_tbl.group_member_id,
			members_tbl.group_id,
			members_tbl.user_id,
			members_tbl.role,
			members_tbl.joined_at,
			coalesce(member_activity_tbl.total_activities, 0),
			coalesce(member_activity_tbl.total_distance, 0),
			coalesce(member_peaks.total_unique_summits, 0),
			coalesce(member_peaks.total_summits, 0)
		FROM members_tbl
		LEFT JOIN member_activity_tbl ON members_tbl.user_id = member_activity_tbl.user_id
		LEFT JOIN member_peaks ON members_tbl.user_id = member_peaks.user_id;
	`
	rows, err := dao.db.Query(sql, groupID, startDate, endDate)
	if err != nil {
		dao.l.Printf("Error getting group members contribution: %v", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		contribution := models.GroupMemberGoalContribution{}
		err = rows.Scan(
			&contribution.GroupMemberID,
			&contribution.GroupID,
			&contribution.UserID,
			&contribution.Role,
			&contribution.JoinedAt,
			&contribution.TotalActivities,
			&contribution.TotalDistance,
			&contribution.TotalUniqueSummits,
			&contribution.TotalSummits,
		)
		if err != nil {
			dao.l.Printf("Error parsing query result: %f", err)
			return nil, err
		}
		groupMembersGoalContribution = append(groupMembersGoalContribution, contribution)
	}
	err = rows.Err()
	if err != nil {
		dao.l.Printf("Error during iteration: %v", err)
		return nil, err
	}

	return groupMembersGoalContribution, nil
}
