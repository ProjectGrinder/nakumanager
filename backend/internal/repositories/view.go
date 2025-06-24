package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/nack098/nakumanager/internal/db"
)

type ViewRepository interface {
	AddGroupByToView(ctx context.Context, data db.AddGroupByToViewParams) error
	AddIssueToView(ctx context.Context, data db.AddIssueToViewParams) error
	CreateView(ctx context.Context, data db.CreateViewParams) error
	DeleteView(ctx context.Context, id string) error
	GetViewByID(ctx context.Context, id string) ([]db.View, error)
	ListGroupByViewID(ctx context.Context, viewID string) ([]string, error)
	ListIssuesByViewID(ctx context.Context, viewID string) ([]db.Issue, error)
	ListViewsByUser(ctx context.Context, userID string) ([]db.View, error)
	RemoveGroupByFromView(ctx context.Context, data db.RemoveGroupByFromViewParams) error
	RemoveIssueFromView(ctx context.Context, data db.RemoveIssueFromViewParams) error
	ListIssuesByGroupFilters(ctx context.Context, teamID string, filters map[string]string) ([]db.Issue, error)
	GetGroupedIssues(
		ctx context.Context,
		teamID string,
		groupBy []string,
	) ([]map[string]interface{}, error)
}

type viewRepo struct {
	db    *db.Queries
	rawDb *sql.DB
}

func NewViewRepository(dbConn *sql.DB) ViewRepository {
	return &viewRepo{
		db:    db.New(dbConn),
		rawDb: dbConn,
	}
}

func (r *viewRepo) AddGroupByToView(ctx context.Context, data db.AddGroupByToViewParams) error {
	return r.db.AddGroupByToView(ctx, data)
}

func (r *viewRepo) AddIssueToView(ctx context.Context, data db.AddIssueToViewParams) error {
	return r.db.AddIssueToView(ctx, data)
}

func (r *viewRepo) CreateView(ctx context.Context, data db.CreateViewParams) error {
	return r.db.CreateView(ctx, data)
}

func (r *viewRepo) DeleteView(ctx context.Context, id string) error {
	return r.db.DeleteView(ctx, id)
}

func (r *viewRepo) GetViewByID(ctx context.Context, viewId string) ([]db.View, error) {
	return r.db.GetViewByID(ctx, viewId)
}

func (r *viewRepo) ListGroupByViewID(ctx context.Context, viewID string) ([]string, error) {
	return r.db.ListGroupByViewID(ctx, viewID)
}

func (r *viewRepo) ListIssuesByViewID(ctx context.Context, viewID string) ([]db.Issue, error) {
	return r.db.ListIssuesByViewID(ctx, viewID)
}

func (r *viewRepo) ListViewsByUser(ctx context.Context, userID string) ([]db.View, error) {
	return r.db.ListViewsByUser(ctx, userID)
}

func (r *viewRepo) RemoveGroupByFromView(ctx context.Context, data db.RemoveGroupByFromViewParams) error {
	return r.db.RemoveGroupByFromView(ctx, data)
}

func (r *viewRepo) RemoveIssueFromView(ctx context.Context, data db.RemoveIssueFromViewParams) error {
	return r.db.RemoveIssueFromView(ctx, data)
}

type GroupedIssue struct {
	Group1     sql.NullString
	Group2     sql.NullString
	IssueCount int
}

func (r *viewRepo) GetGroupedIssues(
	ctx context.Context,
	teamID string,
	groupBy []string,
) ([]map[string]interface{}, error) {
	// Allow only specific fields
	validCols := map[string]bool{
		"status": true, "priority": true, "project_id": true,
		"label": true, "assignee": true, "type": true, "severity": true,
	}

	if len(groupBy) == 0 || len(groupBy) > 2 {
		return nil, fmt.Errorf("groupBy must be 1 or 2 fields")
	}

	for _, col := range groupBy {
		if !validCols[col] {
			return nil, fmt.Errorf("invalid groupBy column: %s", col)
		}
	}

	cols := strings.Join(groupBy, ", ")
	query := fmt.Sprintf(`
		SELECT %s, COUNT(*) as issue_count
		FROM issues
		WHERE team_id = ?
		GROUP BY %s
	`, cols, cols)

	rows, err := r.rawDb.QueryContext(ctx, query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(groupBy)+1)
		for i := range values {
			var v sql.NullString
			values[i] = &v
		}

		if err := rows.Scan(values...); err != nil {
			return nil, err
		}

		row := map[string]interface{}{}
		for i, col := range groupBy {
			row[col] = *(values[i].(*sql.NullString))
		}
		row["count"] = *(values[len(groupBy)].(*sql.NullString)) // แต่จริง ๆ เป็น int

		// convert count เป็น int
		if countVal, ok := row["count"].(sql.NullString); ok && countVal.Valid {
			row["count"], _ = strconv.Atoi(countVal.String)
		}

		results = append(results, row)
	}
	return results, nil
}

func (r *viewRepo) ListIssuesByGroupFilters(ctx context.Context, teamID string, filters map[string]string) ([]db.Issue, error) {
	query := `SELECT * FROM issues WHERE team_id = ?`
	args := []interface{}{teamID}

	for col, val := range filters {
		query += fmt.Sprintf(" AND %s = ?", col)
		args = append(args, val)
	}

	rows, err := r.rawDb.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var issues []db.Issue
	for rows.Next() {
		var issue db.Issue
		if err := rows.Scan(
			&issue.ID,
			&issue.Title,
			&issue.Content,
			&issue.Priority,
			&issue.Status,
			&issue.Assignee,
			&issue.ProjectID,
			&issue.TeamID,
			&issue.StartDate,
			&issue.EndDate,
			&issue.Label,
			&issue.OwnerID,
		); err != nil {
			return nil, err
		}
		issues = append(issues, issue)
	}

	return issues, nil
}
