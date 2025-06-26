-- name: CreateView :exec
INSERT INTO views (id, name, created_by, team_id)
VALUES (?, ?, ?, ?);

-- name: GetViewByID :many
SELECT *
FROM views
WHERE id = ?;


-- name: DeleteView :exec
DELETE FROM views WHERE id = ?;

-- name: ListViewsByUser :many
SELECT id, name, created_by, team_id
FROM views
WHERE created_by = ?
ORDER BY name;

-- name: AddGroupByToView :exec
INSERT INTO view_group_bys (view_id, group_by)
VALUES (?, ?);

-- name: RemoveGroupByFromView :exec
DELETE FROM view_group_bys
WHERE view_id = ?;

-- name: ListGroupByViewID :many
SELECT group_by
FROM view_group_bys
WHERE view_id = ?;

-- name: AddIssueToView :exec
INSERT OR IGNORE INTO view_issues (view_id, issue_id)
VALUES (?, ?);

-- name: RemoveIssueFromView :exec
DELETE FROM view_issues
WHERE view_id = ?;

-- name: ListIssuesByViewID :many
SELECT i.*
FROM issues i
JOIN view_issues vi ON i.id = vi.issue_id
WHERE vi.view_id = ?;

-- name: GetIssuesByStatus :many
SELECT * FROM issues
WHERE team_id = ? AND status = ?;

-- name: GetIssuesByAssignee :many
SELECT i.*
FROM issues i
JOIN issue_assignees ia ON ia.issue_id = i.id
WHERE ia.user_id = ? AND i.team_id = ?;

-- name: GetIssuesByPriority :many
SELECT * FROM issues
WHERE team_id = ? AND priority = ? ;

-- name: GetIssuesByProject :many
SELECT * FROM issues
WHERE team_id = ? AND project_id = ?;

-- name: GetIssuesByLabel :many
SELECT * FROM issues
WHERE team_id = ? AND Label = ?;

-- name: GetIssuesByTeamID :many
SELECT * FROM issues
WHERE team_id = ?;

-- name: GetIssuesByEndDate :many
SELECT * FROM issues
WHERE team_id = ? AND end_date  = ?;

-- name: UpdateViewName :exec
UPDATE views SET name = ? 
WHERE id = ?;

-- name: UpdateViewGroupBy :exec
UPDATE view_group_bys SET group_by = ? 
WHERE view_id = ?;
