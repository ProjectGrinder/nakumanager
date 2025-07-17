package routes

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator"
	models "github.com/nack098/nakumanager/internal/models"
)

var validate = validator.New()

func ToNullString(s *string) sql.NullString {
	if s != nil && strings.TrimSpace(*s) != "" {
		return sql.NullString{String: *s, Valid: true}
	}
	return sql.NullString{}
}

func ToNullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

func nullStringToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func buildUpdateQuery(p models.EditProject) (string, []interface{}) {
	query := "UPDATE projects SET "
	args := []interface{}{}
	sets := []string{}

	if p.Name != nil {
		sets = append(sets, "name = ?")
		args = append(args, *p.Name)
	}
	if p.Status != nil {
		sets = append(sets, "status = ?")
		args = append(args, *p.Status)
	}
	if p.Priority != nil {
		sets = append(sets, "priority = ?")
		args = append(args, *p.Priority)
	}
	if p.StartDate != nil {
		sets = append(sets, "start_date = ?")
		args = append(args, *p.StartDate)
	}
	if p.EndDate != nil {
		sets = append(sets, "end_date = ?")
		args = append(args, *p.EndDate)
	}
	if p.Label != nil {
		sets = append(sets, "label = ?")
		args = append(args, *p.Label)
	}
	if p.LeaderID != nil {
		sets = append(sets, "leader_id = ?")
		args = append(args, *p.LeaderID)
	}

	query += strings.Join(sets, ", ") + " WHERE id = ?"
	args = append(args, p.ID)

	return query, args
}

func buildUpdateIssueQuery(i models.UpdateIssueRequest) (string, []interface{}) {
	query := "UPDATE issues SET "
	args := []interface{}{}
	sets := []string{}

	if i.Title != nil {
		sets = append(sets, "title = ?")
		args = append(args, *i.Title)
	}
	if i.Content != nil {
		sets = append(sets, "content = ?")
		args = append(args, *i.Content)
	}
	if i.Priority != nil {
		sets = append(sets, "priority = ?")
		args = append(args, *i.Priority)
	}
	if i.Status != nil {
		sets = append(sets, "status = ?")
		args = append(args, *i.Status)
	}
	if i.ProjectID != nil {
		sets = append(sets, "project_id = ?")
		args = append(args, *i.ProjectID)
	}
	if i.TeamID != nil {
		sets = append(sets, "team_id = ?")
		args = append(args, *i.TeamID)
	}
	if i.StartDate != nil {
		sets = append(sets, "start_date = ?")
		args = append(args, *i.StartDate)
	}
	if i.EndDate != nil {
		sets = append(sets, "end_date = ?")
		args = append(args, *i.EndDate)
	}
	if i.Label != nil {
		sets = append(sets, "label = ?")
		args = append(args, *i.Label)
	}
	if i.OwnerID != nil {
		sets = append(sets, "owner_id = ?")
		args = append(args, *i.OwnerID)
	}

	if len(sets) == 0 {
		return "", nil
	}

	query += strings.Join(sets, ", ") + " WHERE id = ?"
	args = append(args, i.ID)

	return query, args
}

func BuildGroupByQuery(table string, groupBys []string) (string, error) {
	validCols := map[string]bool{
		"status": true, "priority": true, "project_id": true,
		"label": true, "assignee": true, "team_id": true, "end_date": true,
	}
	if len(groupBys) == 0 {
		return "", fmt.Errorf("no group_by fields provided")
	}
	for _, col := range groupBys {
		if !validCols[col] {
			return "", fmt.Errorf("invalid group_by column: %s", col)
		}
	}

	var conditions []string
	for _, col := range groupBys {
		conditions = append(conditions, fmt.Sprintf("%s IS NOT NULL", col))
	}

	query := fmt.Sprintf(`
		SELECT id
		FROM %s
		WHERE team_id = ? AND %s
	`, table, strings.Join(conditions, " AND "))
	return query, nil
}
