-- name: CreateProject :exec
INSERT INTO projects (name, status, priority, workspace_id, leader_id, start_date, end_date, label)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetProjectByID :one
SELECT id, name, status, priority, workspace_id, leader_id, start_date, end_date, label
FROM projects
WHERE id = ?;

-- name: DeleteProject :exec
DELETE FROM projects WHERE id = ?;

-- name: ListProjectsByWorkspace :many
SELECT id, name, status, priority, workspace_id, leader_id, start_date, end_date, label
FROM projects
WHERE workspace_id = ?
ORDER BY start_date DESC;

-- name: AddMemberToProject :exec
INSERT OR IGNORE INTO project_members (project_id, user_id)
VALUES (?, ?);

-- name: RemoveMemberFromProject :exec
DELETE FROM project_members
WHERE project_id = ? AND user_id = ?;

-- name: ListProjectMembers :many
SELECT u.*
FROM users u
JOIN project_members pm ON u.id = pm.user_id
WHERE pm.project_id = ?;

-- name: UpdateProject :exec
UPDATE projects
SET name = ?, status = ?, priority = ?, workspace_id = ?, leader_id = ?, start_date = ?, end_date = ?, label = ?
WHERE id = ?;
