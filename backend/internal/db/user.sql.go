// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: user.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :exec
INSERT INTO users (id, username, password_hash, email, roles)
VALUES (?, ?, ?, ?, ?)
`

type CreateUserParams struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	Email        string `json:"email"`
	Roles        string `json:"roles"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.ExecContext(ctx, createUser,
		arg.ID,
		arg.Username,
		arg.PasswordHash,
		arg.Email,
		arg.Roles,
	)
	return err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?
`

func (q *Queries) DeleteUser(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getUserByEmailWithPassword = `-- name: GetUserByEmailWithPassword :one
SELECT id, username, email, password_hash, roles
FROM users
WHERE email = ?
`

type GetUserByEmailWithPasswordRow struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Roles        string `json:"roles"`
}

func (q *Queries) GetUserByEmailWithPassword(ctx context.Context, email string) (GetUserByEmailWithPasswordRow, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmailWithPassword, email)
	var i GetUserByEmailWithPasswordRow
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.Roles,
	)
	return i, err
}

const getUserByEmailWithoutPassword = `-- name: GetUserByEmailWithoutPassword :one
SELECT id, username, email, roles
FROM users
WHERE email = ?
`

type GetUserByEmailWithoutPasswordRow struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Roles    string `json:"roles"`
}

func (q *Queries) GetUserByEmailWithoutPassword(ctx context.Context, email string) (GetUserByEmailWithoutPasswordRow, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmailWithoutPassword, email)
	var i GetUserByEmailWithoutPasswordRow
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Roles,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, username, email, roles
FROM users
WHERE id = ?
`

type GetUserByIDRow struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Roles    string `json:"roles"`
}

func (q *Queries) GetUserByID(ctx context.Context, id string) (GetUserByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, id)
	var i GetUserByIDRow
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Roles,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, username, email, roles
FROM users
ORDER BY username
`

type ListUsersRow struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Roles    string `json:"roles"`
}

func (q *Queries) ListUsers(ctx context.Context) ([]ListUsersRow, error) {
	rows, err := q.db.QueryContext(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListUsersRow{}
	for rows.Next() {
		var i ListUsersRow
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Email,
			&i.Roles,
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

const updateEmail = `-- name: UpdateEmail :exec
UPDATE users SET email = ? WHERE id = ?
`

type UpdateEmailParams struct {
	Email string `json:"email"`
	ID    string `json:"id"`
}

func (q *Queries) UpdateEmail(ctx context.Context, arg UpdateEmailParams) error {
	_, err := q.db.ExecContext(ctx, updateEmail, arg.Email, arg.ID)
	return err
}

const updateRoles = `-- name: UpdateRoles :exec
UPDATE users SET roles = ? WHERE id = ?
`

type UpdateRolesParams struct {
	Roles string `json:"roles"`
	ID    string `json:"id"`
}

func (q *Queries) UpdateRoles(ctx context.Context, arg UpdateRolesParams) error {
	_, err := q.db.ExecContext(ctx, updateRoles, arg.Roles, arg.ID)
	return err
}

const updateUsername = `-- name: UpdateUsername :exec
UPDATE users SET username = ? WHERE id = ?
`

type UpdateUsernameParams struct {
	Username string `json:"username"`
	ID       string `json:"id"`
}

func (q *Queries) UpdateUsername(ctx context.Context, arg UpdateUsernameParams) error {
	_, err := q.db.ExecContext(ctx, updateUsername, arg.Username, arg.ID)
	return err
}
