-- name: CreateWorkspace :exec
INSERT INTO workspaces (id, name, owner_id) 
VALUES (?, ?, ?);

-- name: GetWorkspaceByID :one
SELECT * FROM workspaces WHERE id = ?;

-- name: GetWorkspaceByUserID :many
SELECT w.*
FROM workspaces w
WHERE w.owner_id = ?;

-- name: DeleteWorkspace :exec
DELETE FROM workspaces WHERE id = ?;

-- name: ListWorkspaceMembers :many
SELECT u.*
FROM users u
JOIN workspace_members wm ON u.id = wm.user_id
WHERE wm.workspace_id = ?;

-- name: AddMemberToWorkspace :exec
INSERT OR IGNORE INTO workspace_members (workspace_id, user_id)
VALUES (?, ?);

-- name: RemoveMemberFromWorkspace :exec
DELETE FROM workspace_members
WHERE workspace_id = ? AND user_id = ?;

-- name: RenameWorkspace :exec
UPDATE workspaces
SET name = ?
WHERE id = ?;

-- name: ListWorkspacesWithMembersByUserID :many
SELECT w.id, w.name, w.owner_id, wm.user_id
FROM workspaces w
LEFT JOIN workspace_members wm ON w.id = wm.workspace_id
WHERE w.owner_id = ? OR wm.user_id = ?
ORDER BY w.id;
