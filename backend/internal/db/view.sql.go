// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: view.sql

package db

import (
	"context"
	"database/sql"
)

const addGroupByToView = `-- name: AddGroupByToView :exec
INSERT INTO view_group_bys (view_id, group_by)
VALUES (?, ?)
`

type AddGroupByToViewParams struct {
	ViewID  string `json:"view_id"`
	GroupBy string `json:"group_by"`
}

func (q *Queries) AddGroupByToView(ctx context.Context, arg AddGroupByToViewParams) error {
	_, err := q.db.ExecContext(ctx, addGroupByToView, arg.ViewID, arg.GroupBy)
	return err
}

const addIssueToView = `-- name: AddIssueToView :exec
INSERT OR IGNORE INTO view_issues (view_id, issue_id)
VALUES (?, ?)
`

type AddIssueToViewParams struct {
	ViewID  string `json:"view_id"`
	IssueID string `json:"issue_id"`
}

func (q *Queries) AddIssueToView(ctx context.Context, arg AddIssueToViewParams) error {
	_, err := q.db.ExecContext(ctx, addIssueToView, arg.ViewID, arg.IssueID)
	return err
}

const createView = `-- name: CreateView :exec
INSERT INTO views (id, name, created_by, team_id)
VALUES (?, ?, ?, ?)
`

type CreateViewParams struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedBy string `json:"created_by"`
	TeamID    string `json:"team_id"`
}

func (q *Queries) CreateView(ctx context.Context, arg CreateViewParams) error {
	_, err := q.db.ExecContext(ctx, createView,
		arg.ID,
		arg.Name,
		arg.CreatedBy,
		arg.TeamID,
	)
	return err
}

const deleteView = `-- name: DeleteView :exec
DELETE FROM views WHERE id = ?
`

func (q *Queries) DeleteView(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteView, id)
	return err
}

const getIssuesByAssignee = `-- name: GetIssuesByAssignee :many
SELECT i.id, i.title, i.content, i.priority, i.status, i.project_id, i.team_id, i.start_date, i.end_date, i.label, i.owner_id
FROM issues i
JOIN issue_assignees ia ON ia.issue_id = i.id
WHERE ia.user_id = ? AND i.team_id = ?
`

type GetIssuesByAssigneeParams struct {
	UserID string `json:"user_id"`
	TeamID string `json:"team_id"`
}

func (q *Queries) GetIssuesByAssignee(ctx context.Context, arg GetIssuesByAssigneeParams) ([]Issue, error) {
	rows, err := q.db.QueryContext(ctx, getIssuesByAssignee, arg.UserID, arg.TeamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Issue{}
	for rows.Next() {
		var i Issue
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Content,
			&i.Priority,
			&i.Status,
			&i.ProjectID,
			&i.TeamID,
			&i.StartDate,
			&i.EndDate,
			&i.Label,
			&i.OwnerID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getIssuesByEndDate = `-- name: GetIssuesByEndDate :many
SELECT id, title, content, priority, status, project_id, team_id, start_date, end_date, label, owner_id FROM issues
WHERE team_id = ? AND end_date  = ?
`

type GetIssuesByEndDateParams struct {
	TeamID  string       `json:"team_id"`
	EndDate sql.NullTime `json:"end_date"`
}

func (q *Queries) GetIssuesByEndDate(ctx context.Context, arg GetIssuesByEndDateParams) ([]Issue, error) {
	rows, err := q.db.QueryContext(ctx, getIssuesByEndDate, arg.TeamID, arg.EndDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Issue{}
	for rows.Next() {
		var i Issue
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Content,
			&i.Priority,
			&i.Status,
			&i.ProjectID,
			&i.TeamID,
			&i.StartDate,
			&i.EndDate,
			&i.Label,
			&i.OwnerID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getIssuesByLabel = `-- name: GetIssuesByLabel :many
SELECT id, title, content, priority, status, project_id, team_id, start_date, end_date, label, owner_id FROM issues
WHERE team_id = ? AND Label = ?
`

type GetIssuesByLabelParams struct {
	TeamID string         `json:"team_id"`
	Label  sql.NullString `json:"label"`
}

func (q *Queries) GetIssuesByLabel(ctx context.Context, arg GetIssuesByLabelParams) ([]Issue, error) {
	rows, err := q.db.QueryContext(ctx, getIssuesByLabel, arg.TeamID, arg.Label)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Issue{}
	for rows.Next() {
		var i Issue
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Content,
			&i.Priority,
			&i.Status,
			&i.ProjectID,
			&i.TeamID,
			&i.StartDate,
			&i.EndDate,
			&i.Label,
			&i.OwnerID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getIssuesByPriority = `-- name: GetIssuesByPriority :many
SELECT id, title, content, priority, status, project_id, team_id, start_date, end_date, label, owner_id FROM issues
WHERE team_id = ? AND priority = ?
`

type GetIssuesByPriorityParams struct {
	TeamID   string         `json:"team_id"`
	Priority sql.NullString `json:"priority"`
}

func (q *Queries) GetIssuesByPriority(ctx context.Context, arg GetIssuesByPriorityParams) ([]Issue, error) {
	rows, err := q.db.QueryContext(ctx, getIssuesByPriority, arg.TeamID, arg.Priority)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Issue{}
	for rows.Next() {
		var i Issue
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Content,
			&i.Priority,
			&i.Status,
			&i.ProjectID,
			&i.TeamID,
			&i.StartDate,
			&i.EndDate,
			&i.Label,
			&i.OwnerID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getIssuesByProject = `-- name: GetIssuesByProject :many
;

SELECT id, title, content, priority, status, project_id, team_id, start_date, end_date, label, owner_id FROM issues
WHERE team_id = ? AND project_id = ?
`

type GetIssuesByProjectParams struct {
	TeamID    string         `json:"team_id"`
	ProjectID sql.NullString `json:"project_id"`
}

func (q *Queries) GetIssuesByProject(ctx context.Context, arg GetIssuesByProjectParams) ([]Issue, error) {
	rows, err := q.db.QueryContext(ctx, getIssuesByProject, arg.TeamID, arg.ProjectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Issue{}
	for rows.Next() {
		var i Issue
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Content,
			&i.Priority,
			&i.Status,
			&i.ProjectID,
			&i.TeamID,
			&i.StartDate,
			&i.EndDate,
			&i.Label,
			&i.OwnerID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getIssuesByStatus = `-- name: GetIssuesByStatus :many
SELECT id, title, content, priority, status, project_id, team_id, start_date, end_date, label, owner_id FROM issues
WHERE team_id = ? AND status = ?
`

type GetIssuesByStatusParams struct {
	TeamID string `json:"team_id"`
	Status string `json:"status"`
}

func (q *Queries) GetIssuesByStatus(ctx context.Context, arg GetIssuesByStatusParams) ([]Issue, error) {
	rows, err := q.db.QueryContext(ctx, getIssuesByStatus, arg.TeamID, arg.Status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Issue{}
	for rows.Next() {
		var i Issue
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Content,
			&i.Priority,
			&i.Status,
			&i.ProjectID,
			&i.TeamID,
			&i.StartDate,
			&i.EndDate,
			&i.Label,
			&i.OwnerID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getIssuesByTeamID = `-- name: GetIssuesByTeamID :many
SELECT id, title, content, priority, status, project_id, team_id, start_date, end_date, label, owner_id FROM issues
WHERE team_id = ?
`

func (q *Queries) GetIssuesByTeamID(ctx context.Context, teamID string) ([]Issue, error) {
	rows, err := q.db.QueryContext(ctx, getIssuesByTeamID, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Issue{}
	for rows.Next() {
		var i Issue
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Content,
			&i.Priority,
			&i.Status,
			&i.ProjectID,
			&i.TeamID,
			&i.StartDate,
			&i.EndDate,
			&i.Label,
			&i.OwnerID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTeamIDByViewID = `-- name: GetTeamIDByViewID :one
SELECT team_id
FROM views
WHERE id = ?
`

func (q *Queries) GetTeamIDByViewID(ctx context.Context, id string) (string, error) {
	row := q.db.QueryRowContext(ctx, getTeamIDByViewID, id)
	var team_id string
	err := row.Scan(&team_id)
	return team_id, err
}

const getViewByID = `-- name: GetViewByID :many
SELECT id, name, created_by, team_id
FROM views
WHERE id = ?
`

func (q *Queries) GetViewByID(ctx context.Context, id string) ([]View, error) {
	rows, err := q.db.QueryContext(ctx, getViewByID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []View{}
	for rows.Next() {
		var i View
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedBy,
			&i.TeamID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listGroupByViewID = `-- name: ListGroupByViewID :many
SELECT group_by
FROM view_group_bys
WHERE view_id = ?
`

func (q *Queries) ListGroupByViewID(ctx context.Context, viewID string) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, listGroupByViewID, viewID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []string{}
	for rows.Next() {
		var group_by string
		if err := rows.Scan(&group_by); err != nil {
			return nil, err
		}
		items = append(items, group_by)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listIssuesByViewID = `-- name: ListIssuesByViewID :many
SELECT i.id, i.title, i.content, i.priority, i.status, i.project_id, i.team_id, i.start_date, i.end_date, i.label, i.owner_id
FROM issues i
JOIN view_issues vi ON i.id = vi.issue_id
WHERE vi.view_id = ?
`

func (q *Queries) ListIssuesByViewID(ctx context.Context, viewID string) ([]Issue, error) {
	rows, err := q.db.QueryContext(ctx, listIssuesByViewID, viewID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Issue{}
	for rows.Next() {
		var i Issue
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Content,
			&i.Priority,
			&i.Status,
			&i.ProjectID,
			&i.TeamID,
			&i.StartDate,
			&i.EndDate,
			&i.Label,
			&i.OwnerID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listViewByTeamID = `-- name: ListViewByTeamID :many
SELECT id, name, created_by, team_id
FROM views
WHERE team_id = ?
`

func (q *Queries) ListViewByTeamID(ctx context.Context, teamID string) ([]View, error) {
	rows, err := q.db.QueryContext(ctx, listViewByTeamID, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []View{}
	for rows.Next() {
		var i View
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedBy,
			&i.TeamID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listViewsByUser = `-- name: ListViewsByUser :many
SELECT id, name, created_by, team_id
FROM views
WHERE created_by = ?
ORDER BY name
`

func (q *Queries) ListViewsByUser(ctx context.Context, createdBy string) ([]View, error) {
	rows, err := q.db.QueryContext(ctx, listViewsByUser, createdBy)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []View{}
	for rows.Next() {
		var i View
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.CreatedBy,
			&i.TeamID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const removeGroupByFromView = `-- name: RemoveGroupByFromView :exec
DELETE FROM view_group_bys
WHERE view_id = ?
`

func (q *Queries) RemoveGroupByFromView(ctx context.Context, viewID string) error {
	_, err := q.db.ExecContext(ctx, removeGroupByFromView, viewID)
	return err
}

const removeIssueFromView = `-- name: RemoveIssueFromView :exec
DELETE FROM view_issues
WHERE view_id = ?
`

func (q *Queries) RemoveIssueFromView(ctx context.Context, viewID string) error {
	_, err := q.db.ExecContext(ctx, removeIssueFromView, viewID)
	return err
}

const updateViewGroupBy = `-- name: UpdateViewGroupBy :exec
UPDATE view_group_bys SET group_by = ? 
WHERE view_id = ?
`

type UpdateViewGroupByParams struct {
	GroupBy string `json:"group_by"`
	ViewID  string `json:"view_id"`
}

func (q *Queries) UpdateViewGroupBy(ctx context.Context, arg UpdateViewGroupByParams) error {
	_, err := q.db.ExecContext(ctx, updateViewGroupBy, arg.GroupBy, arg.ViewID)
	return err
}

const updateViewName = `-- name: UpdateViewName :exec
UPDATE views SET name = ? 
WHERE id = ?
`

type UpdateViewNameParams struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func (q *Queries) UpdateViewName(ctx context.Context, arg UpdateViewNameParams) error {
	_, err := q.db.ExecContext(ctx, updateViewName, arg.Name, arg.ID)
	return err
}

const updateViewTeamID = `-- name: UpdateViewTeamID :exec
UPDATE views SET team_id = ? 
WHERE id = ?
`

type UpdateViewTeamIDParams struct {
	TeamID string `json:"team_id"`
	ID     string `json:"id"`
}

func (q *Queries) UpdateViewTeamID(ctx context.Context, arg UpdateViewTeamIDParams) error {
	_, err := q.db.ExecContext(ctx, updateViewTeamID, arg.TeamID, arg.ID)
	return err
}
