-- name: CreateProject :exec
INSERT INTO projects (
  id, name, status, priority, workspace_id, team_id, leader_id, start_date, end_date, label, created_by
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);


-- name: GetProjectByID :one
SELECT id, name, status, priority, workspace_id, team_id, leader_id, start_date, end_date, label, created_by
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


-- name: GetProjectsByUserID :many
SELECT p.*
FROM projects p
JOIN project_members pm ON p.id = pm.project_id
WHERE pm.user_id = ?;
