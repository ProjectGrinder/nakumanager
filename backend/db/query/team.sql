-- name: GetTeamsByUserID :many
SELECT t.id, t.name, t.workspace_id, t.leader_id
FROM teams t
JOIN team_members tm ON t.id = tm.team_id
WHERE tm.user_id = ?;

-- name: DeleteTeam :exec
DELETE FROM teams WHERE id = ?;

-- name: ListTeams :many
SELECT id, name, workspace_id, leader_id
FROM teams
ORDER BY name;

-- name: AddMemberToTeam :exec
INSERT OR IGNORE INTO team_members (team_id, user_id)
VALUES (?, ?);

-- name: RemoveMemberFromTeam :exec
DELETE FROM team_members
WHERE team_id = ? AND user_id = ?;

-- name: ListTeamMembers :many
SELECT u.id, u.username, u.email
FROM users u
JOIN team_members tm ON u.id = tm.user_id
WHERE tm.team_id = ?;

-- name: GetOwnerByTeamID :one
SELECT w.owner_id
FROM teams t
JOIN workspaces w ON t.workspace_id = w.id
WHERE t.id = ?;


-- name: GetLeaderByTeamID :one
SELECT leader_id
FROM teams
WHERE id = ?;

-- name: IsMemberInTeam :one
SELECT COUNT(*) AS count
FROM team_members
WHERE team_id = ? AND user_id = ?;

-- name: IsTeamExists :one
SELECT COUNT(*) AS count
FROM teams
WHERE id = ?;

-- name: ListIssuesByUserID :many
SELECT i.id, i.title, i.status, i.priority, i.project_id
FROM issues i
JOIN project_members pm ON i.project_id = pm.project_id
WHERE pm.user_id = ?;

-- name: RenameTeam :exec
UPDATE teams
SET name = ?
WHERE id = ?;

-- name: CreateTeam :exec
INSERT INTO teams (id, name, workspace_id)
VALUES (?, ?, ?);

-- name: GetTeamByID :one
SELECT id, name, workspace_id, leader_id
FROM teams
WHERE id = ?;

-- name: SetLeaderToTeam :exec
UPDATE teams
SET leader_id = ?
WHERE id = ?;

