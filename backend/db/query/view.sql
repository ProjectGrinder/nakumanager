-- name: CreateView :exec
INSERT INTO views (name, user_id, team_id)
VALUES (?, ?, ?);

-- name: GetViewByID :one
SELECT id, name, user_id, team_id
FROM views
WHERE id = ?;

-- name: DeleteView :exec
DELETE FROM views WHERE id = ?;

-- name: ListViewsByUser :many
SELECT id, name, user_id, team_id
FROM views
WHERE user_id = ?
ORDER BY name;

-- name: AddGroupByToView :exec
INSERT INTO view_group_bys (view_id, group_by)
VALUES (?, ?);

-- name: RemoveGroupByFromView :exec
DELETE FROM view_group_bys
WHERE view_id = ? AND group_by = ?;

-- name: ListGroupBysByViewID :many
SELECT group_by
FROM view_group_bys
WHERE view_id = ?;

-- name: AddIssueToView :exec
INSERT OR IGNORE INTO view_issues (view_id, issue_id)
VALUES (?, ?);

-- name: RemoveIssueFromView :exec
DELETE FROM view_issues
WHERE view_id = ? AND issue_id = ?;

-- name: ListIssuesByViewID :many
SELECT i.*
FROM issues i
JOIN view_issues vi ON i.id = vi.issue_id
WHERE vi.view_id = ?;
