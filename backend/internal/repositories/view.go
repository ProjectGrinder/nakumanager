package repositories

import (
	"context"
	"database/sql"
	"strings"

	"github.com/nack098/nakumanager/internal/db"
)

type ViewRepository interface {
	AddGroupByToView(ctx context.Context, data db.AddGroupByToViewParams) error
	AddIssueToView(ctx context.Context, data db.AddIssueToViewParams) error
	AddIssueToViewTx(ctx context.Context, tx *sql.Tx, params db.AddIssueToViewParams) error
	CreateView(ctx context.Context, data db.CreateViewParams) error
	DeleteView(ctx context.Context, id string) error
	GetViewByID(ctx context.Context, id string) ([]db.View, error)
	ListGroupByViewID(ctx context.Context, viewID string) ([]string, error)
	ListIssuesByViewID(ctx context.Context, viewID string) ([]db.Issue, error)
	ListViewsByUser(ctx context.Context, userID string) ([]db.View, error)
	UpdateViewName(ctx context.Context, id string, name string) error
	RemoveGroupByFromView(ctx context.Context, id string) error
	RemoveIssueFromView(ctx context.Context, id string) error
	ListViewByTeamID(ctx context.Context, teamID string) ([]db.View, error)
	GetViewsByGroupBys(ctx context.Context, groupBys []string) ([]db.View, error)
	UpdateViewTeamID(ctx context.Context, id string, teamID string) error
	GetTeamIDByViewID(ctx context.Context, id string) (string, error)
}

type viewRepo struct {
	db    *db.Queries
	rawDb *sql.DB
}

func NewViewRepository(dbConn *sql.DB) ViewRepository {
	return &viewRepo{
		db:    db.New(dbConn),
		rawDb: dbConn,
	}
}

func (r *viewRepo) AddGroupByToView(ctx context.Context, data db.AddGroupByToViewParams) error {
	return r.db.AddGroupByToView(ctx, data)
}

func (r *viewRepo) AddIssueToView(ctx context.Context, data db.AddIssueToViewParams) error {
	return r.db.AddIssueToView(ctx, data)
}

func (r *viewRepo) AddIssueToViewTx(ctx context.Context, tx *sql.Tx, params db.AddIssueToViewParams) error {
	q := db.New(tx) // ใช้ sqlc กับ transaction นี้
	return q.AddIssueToView(ctx, params)
}

func (r *viewRepo) CreateView(ctx context.Context, data db.CreateViewParams) error {
	return r.db.CreateView(ctx, data)
}

func (r *viewRepo) DeleteView(ctx context.Context, id string) error {
	return r.db.DeleteView(ctx, id)
}

func (r *viewRepo) GetViewByID(ctx context.Context, viewId string) ([]db.View, error) {
	return r.db.GetViewByID(ctx, viewId)
}

func (r *viewRepo) ListGroupByViewID(ctx context.Context, viewID string) ([]string, error) {
	return r.db.ListGroupByViewID(ctx, viewID)
}

func (r *viewRepo) ListIssuesByViewID(ctx context.Context, viewID string) ([]db.Issue, error) {
	return r.db.ListIssuesByViewID(ctx, viewID)
}

func (r *viewRepo) ListViewsByUser(ctx context.Context, userID string) ([]db.View, error) {
	return r.db.ListViewsByUser(ctx, userID)
}

func (r *viewRepo) UpdateViewName(ctx context.Context, id string, name string) error {
	return r.db.UpdateViewName(ctx, db.UpdateViewNameParams{
		ID:   id,
		Name: name,
	})
}

func (r *viewRepo) RemoveGroupByFromView(ctx context.Context, id string) error {
	return r.db.RemoveGroupByFromView(ctx, id)
}

func (r *viewRepo) RemoveIssueFromView(ctx context.Context, id string) error {
	return r.db.RemoveIssueFromView(ctx, id)
}

func (r *viewRepo) ListViewByTeamID(ctx context.Context, teamID string) ([]db.View, error) {
	return r.db.ListViewByTeamID(ctx, teamID)
}

func (r *viewRepo) GetViewsByGroupBys(ctx context.Context, groupBys []string) ([]db.View, error) {
	query := `
    SELECT v.id, v.name, v.created_by, v.team_id
	FROM views v
	JOIN view_group_bys vg ON v.id = vg.view_id
	WHERE vg.group_by IN (?` + strings.Repeat(",?", len(groupBys)-1) + `)
	GROUP BY v.id
	HAVING COUNT(DISTINCT vg.group_by) = ?;
    `

	args := make([]interface{}, len(groupBys)+1)
	for i, g := range groupBys {
		args[i] = g
	}
	args[len(groupBys)] = len(groupBys)

	rows, err := r.rawDb.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var views []db.View
	for rows.Next() {
		var v db.View
		if err := rows.Scan(&v.ID, &v.Name, &v.CreatedBy, &v.TeamID); err != nil {
			return nil, err
		}
		views = append(views, v)
	}
	return views, nil
}

func (r *viewRepo) UpdateViewTeamID(ctx context.Context, id string, teamID string) error {
	return r.db.UpdateViewTeamID(ctx, db.UpdateViewTeamIDParams{
		ID:     id,
		TeamID: teamID,
	})
}

func (r *viewRepo) GetTeamIDByViewID(ctx context.Context, id string) (string, error) {
	return r.db.GetTeamIDByViewID(ctx, id)
}
