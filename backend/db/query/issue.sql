-- name: CreateIssue :exec
INSERT INTO issues (
    id, title, content, priority, status, project_id, team_id, assignee,
    start_date, end_date, label, owner_id
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetIssueByID :one
SELECT *
FROM issues
WHERE id = ?;

-- name: DeleteIssue :exec
DELETE FROM issues WHERE id = ?;

-- name: ListIssuesByProjectID :many
SELECT *
FROM issues
WHERE project_id = ?
ORDER BY start_date DESC;

-- name: ListIssuesByTeamID :many
SELECT *
FROM issues
WHERE team_id = ?
ORDER BY start_date DESC;

-- name: AddAssigneeToIssue :exec
INSERT OR IGNORE INTO issue_assignees (issue_id, user_id)
VALUES (?, ?);

-- name: RemoveAssigneeFromIssue :exec
DELETE FROM issue_assignees
WHERE issue_id = ? AND user_id = ?;

-- name: ListAssigneesByIssueID :many
SELECT u.*
FROM users u
JOIN issue_assignees ia ON u.id = ia.user_id
WHERE ia.issue_id = ?;


-- name: GetIssueByUserID :many
SELECT i.*
FROM issues i
JOIN issue_assignees ia ON i.id = ia.issue_id
WHERE ia.user_id = ?;

-- name: UpdateIssue :exec
UPDATE issues
SET
    title = ?,
    content = ?,
    priority = ?,
    status = ?,
    project_id = ?,
    team_id = ?,
    start_date = ?,
    end_date = ?,
    label = ?,
    owner_id = ?
WHERE id = ?;
