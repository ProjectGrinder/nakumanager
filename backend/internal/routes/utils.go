package routes

import (
	"database/sql"
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

func toNullTime(s *string) sql.NullTime {
	if s != nil {
		if t, err := time.Parse(time.RFC3339, *s); err == nil {
			return sql.NullTime{Time: t, Valid: true}
		}
	}
	return sql.NullTime{}
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
