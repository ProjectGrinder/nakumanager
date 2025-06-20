-- name: CreateTeam :exec
INSERT INTO teams (id, name, leader_id)
VALUES (?, ?, ?);

-- name: GetTeamByID :one
SELECT id, name, leader_id
FROM teams
WHERE id = ?;

-- name: DeleteTeam :exec
DELETE FROM teams WHERE id = ?;

-- name: ListTeams :many
SELECT id, name, leader_id
FROM teams
ORDER BY name;

-- name: AddMemberToTeam :exec
INSERT OR IGNORE INTO team_members (team_id, user_id)
VALUES (?, ?);

-- name: RemoveMemberFromTeam :exec
DELETE FROM team_members
WHERE team_id = ? AND user_id = ?;

-- name: ListTeamMembers :many
SELECT u.*
FROM users u
JOIN team_members tm ON u.id = tm.user_id
WHERE tm.team_id = ?;
