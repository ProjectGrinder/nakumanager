package routes_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/nack098/nakumanager/internal/db"
	models "github.com/nack098/nakumanager/internal/models"
	"github.com/nack098/nakumanager/internal/routes"
	mocks "github.com/nack098/nakumanager/internal/routes/mock_repo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewViewHandler(t *testing.T) {
	db := &sql.DB{}
	mockRepo := new(mocks.MockViewRepo)
	handler := routes.NewViewHandler(db, mockRepo)
	assert.NotNil(t, handler)
	assert.Equal(t, mockRepo, handler.Repo)
}

func TestCreateView(t *testing.T) {
	app := fiber.New()
	mockRepo := new(mocks.MockViewRepo)
	dbMock, sqlMock, err := sqlmock.New()
	require.NoError(t, err)

	handler := &routes.ViewHandler{
		DB:   dbMock,
		Repo: mockRepo,
	}

	app.Use(withUserID("user-123"))
	app.Post("/views", handler.CreateView)

	tests := []struct {
		name       string
		payload    any
		setupMocks func()
		wantStatus int
	}{
		{
			name: "success",
			payload: models.CreateView{
				Name:     "Test View",
				TeamID:   "team-1",
				GroupBys: []string{"status", "priority"},
			},
			setupMocks: func() {
				mockRepo.On("CreateView", mock.Anything, mock.Anything).Return(nil)
				mockRepo.On("AddGroupByToView", mock.Anything, mock.Anything).Return(nil).Twice()
				sqlMock.ExpectQuery("SELECT .* FROM issues .*").
					WithArgs("team-1").
					WillReturnRows(sqlmock.NewRows([]string{"issue_id"}).AddRow("i1").AddRow("i2"))
				mockRepo.On("AddIssueToView", mock.Anything, mock.Anything).Return(nil).Twice()
			},
			wantStatus: fiber.StatusCreated,
		},
		{
			name:       "bad body",
			payload:    "not-a-json",
			setupMocks: func() {},
			wantStatus: fiber.StatusBadRequest,
		},
		{
			name:       "missing name/teamID",
			payload:    models.CreateView{},
			setupMocks: func() {},
			wantStatus: fiber.StatusBadRequest,
		},
		{
			name:    "create view fail",
			payload: models.CreateView{Name: "View", TeamID: "team"},
			setupMocks: func() {
				mockRepo.On("CreateView", mock.Anything, mock.Anything).
					Return(errors.New("create error"))
			},
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name:    "invalid group_by",
			payload: models.CreateView{Name: "View", TeamID: "team", GroupBys: []string{"invalid"}},
			setupMocks: func() {
				mockRepo.On("CreateView", mock.Anything, mock.Anything).Return(nil)
			},
			wantStatus: fiber.StatusBadRequest,
		},
		{
			name:    "group_by insert fail",
			payload: models.CreateView{Name: "View", TeamID: "team", GroupBys: []string{"status"}},
			setupMocks: func() {
				mockRepo.On("CreateView", mock.Anything, mock.Anything).Return(nil)
				mockRepo.On("AddGroupByToView", mock.Anything, mock.Anything).
					Return(errors.New("groupBy error"))
			},
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name:    "query fail",
			payload: models.CreateView{Name: "View", TeamID: "team", GroupBys: []string{"status"}},
			setupMocks: func() {
				mockRepo.On("CreateView", mock.Anything, mock.Anything).Return(nil)
				mockRepo.On("AddGroupByToView", mock.Anything, mock.Anything).Return(nil)
				sqlMock.ExpectQuery("SELECT .* FROM issues .*").
					WithArgs("team").WillReturnError(errors.New("query error"))
			},
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name:    "row scan fail",
			payload: models.CreateView{Name: "View", TeamID: "team", GroupBys: []string{"status"}},
			setupMocks: func() {
				mockRepo.On("CreateView", mock.Anything, mock.Anything).Return(nil)
				mockRepo.On("AddGroupByToView", mock.Anything, mock.Anything).Return(nil)
				sqlMock.ExpectQuery("SELECT .* FROM issues .*").
					WithArgs("team").
					WillReturnRows(sqlmock.NewRows([]string{"issue_id"}).
						AddRow(nil))
			},
			wantStatus: fiber.StatusCreated,
		},
		{
			name:    "rows.Err non-nil",
			payload: models.CreateView{Name: "View", TeamID: "team", GroupBys: []string{"status"}},
			setupMocks: func() {
				mockRepo.On("CreateView", mock.Anything, mock.Anything).Return(nil)
				mockRepo.On("AddGroupByToView", mock.Anything, mock.Anything).Return(nil)
				rows := sqlmock.NewRows([]string{"issue_id"}).AddRow("i1")
				rows.RowError(0, errors.New("row error"))
				sqlMock.ExpectQuery("SELECT .* FROM issues .*").
					WithArgs("team").WillReturnRows(rows)
			},
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name:    "add issue fail",
			payload: models.CreateView{Name: "View", TeamID: "team", GroupBys: []string{"status"}},
			setupMocks: func() {
				mockRepo.On("CreateView", mock.Anything, mock.Anything).Return(nil)
				mockRepo.On("AddGroupByToView", mock.Anything, mock.Anything).Return(nil)
				sqlMock.ExpectQuery("SELECT .* FROM issues .*").
					WithArgs("team").
					WillReturnRows(sqlmock.NewRows([]string{"issue_id"}).AddRow("i1"))
				mockRepo.On("AddIssueToView", mock.Anything, mock.Anything).
					Return(errors.New("add issue fail"))
			},
			wantStatus: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			sqlMock.ExpectationsWereMet()
			tt.setupMocks()

			var body []byte
			if str, ok := tt.payload.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.payload)
			}

			req := httptest.NewRequest(http.MethodPost, "/views", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestGetViewsByGroupBy(t *testing.T) {
	app := fiber.New()
	mockRepo := new(mocks.MockViewRepo)

	handler := &routes.ViewHandler{
		DB:   nil,
		Repo: mockRepo,
	}

	app.Use(func(c *fiber.Ctx) error {
		if c.Path() != "/unauthorized" {
			c.Locals("userID", "user-123")
		}
		return c.Next()
	})
	app.Post("/views/group-by", handler.GetViewsByGroupBy)
	app.Post("/unauthorized", handler.GetViewsByGroupBy)

	tests := []struct {
		name       string
		path       string
		body       interface{}
		setupMocks func()
		wantStatus int
	}{
		{
			name:       "invalid body",
			path:       "/views/group-by",
			body:       "not-json",
			setupMocks: func() {},
			wantStatus: fiber.StatusBadRequest,
		},
		{
			name: "missing group_bys",
			path: "/views/group-by",
			body: models.ViewGroupBy{
				GroupBys: []string{},
			},
			setupMocks: func() {},
			wantStatus: fiber.StatusBadRequest,
		},
		{
			name: "repo error",
			path: "/views/group-by",
			body: models.ViewGroupBy{
				GroupBys: []string{"status"},
			},
			setupMocks: func() {
				mockRepo.On("GetViewsByGroupBys", mock.Anything, []string{"status"}).
					Return(nil, errors.New("repo failed"))
			},
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name: "success",
			path: "/views/group-by",
			body: models.ViewGroupBy{
				GroupBys: []string{"priority", "project_id"},
			},
			setupMocks: func() {
				mockRepo.On("GetViewsByGroupBys", mock.Anything, []string{"priority", "project_id"}).
					Return([]db.View{
						{ID: "v1", Name: "View 1"},
					}, nil)
			},
			wantStatus: fiber.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.setupMocks()

			var body []byte
			if str, ok := tt.body.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(http.MethodPost, tt.path, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestGetViewByTeamID(t *testing.T) {
	app := fiber.New()
	mockRepo := new(mocks.MockViewRepo)

	handler := &routes.ViewHandler{
		DB:   nil,
		Repo: mockRepo,
	}

	app.Get("/views/team/:id", handler.GetViewByTeamID)

	tests := []struct {
		name       string
		teamID     string
		setupMocks func()
		wantStatus int
	}{
		{
			name:       "teamID is empty",
			teamID:     "",
			setupMocks: func() {},
			wantStatus: fiber.StatusNotFound,
		},
		{
			name:       "teamID is 'undefined'",
			teamID:     "undefined",
			setupMocks: func() {},
			wantStatus: fiber.StatusBadRequest,
		},
		{
			name:   "repo error",
			teamID: "team-x",
			setupMocks: func() {
				mockRepo.On("ListViewByTeamID", mock.Anything, "team-x").
					Return(nil, errors.New("repo fail"))
			},
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name:   "success",
			teamID: "team-1",
			setupMocks: func() {
				mockRepo.On("ListViewByTeamID", mock.Anything, "team-1").
					Return([]db.View{
						{ID: "v1", Name: "View One", TeamID: "team-1"},
					}, nil)
			},
			wantStatus: fiber.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.setupMocks()

			url := "/views/team/" + tt.teamID
			req := httptest.NewRequest(http.MethodGet, url, nil)
			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestDeleteView(t *testing.T) {
	app := fiber.New()
	mockRepo := new(mocks.MockViewRepo)

	handler := &routes.ViewHandler{
		DB:   nil,
		Repo: mockRepo,
	}

	app.Delete("/views/:id", handler.DeleteView)

	tests := []struct {
		name       string
		viewID     string
		setupMocks func()
		wantStatus int
	}{
		{
			name:       "missing viewID",
			viewID:     "",
			setupMocks: func() {},
			wantStatus: fiber.StatusNotFound,
		},
		{
			name:       "viewID is 'undefined'",
			viewID:     "undefined",
			setupMocks: func() {},
			wantStatus: fiber.StatusBadRequest,
		},
		{
			name:   "repo.GetViewByID fails",
			viewID: "v123",
			setupMocks: func() {
				mockRepo.On("GetViewByID", mock.Anything, "v123").
					Return(nil, errors.New("DB error"))
			},
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name:   "view not found",
			viewID: "v124",
			setupMocks: func() {
				mockRepo.On("GetViewByID", mock.Anything, "v124").
					Return([]db.View{}, nil)
			},
			wantStatus: fiber.StatusNotFound,
		},
		{
			name:   "repo.DeleteView fails",
			viewID: "v125",
			setupMocks: func() {
				mockRepo.On("GetViewByID", mock.Anything, "v125").
					Return([]db.View{{ID: "v125", Name: "Demo"}}, nil)
				mockRepo.On("DeleteView", mock.Anything, "v125").
					Return(errors.New("delete fail"))
			},
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name:   "success",
			viewID: "v126",
			setupMocks: func() {
				mockRepo.On("GetViewByID", mock.Anything, "v126").
					Return([]db.View{{ID: "v126", Name: "Demo"}}, nil)
				mockRepo.On("DeleteView", mock.Anything, "v126").
					Return(nil)
			},
			wantStatus: fiber.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.setupMocks()

			url := "/views/" + tt.viewID
			req := httptest.NewRequest(http.MethodDelete, url, nil)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}

func TestUpdateView(t *testing.T) {
	app := fiber.New()
	mockRepo := new(mocks.MockViewRepo)
	mockDB, sqlMock, err := sqlmock.New()
	require.NoError(t, err)

	handler := &routes.ViewHandler{
		DB:   mockDB,
		Repo: mockRepo,
	}

	app.Use(withUserID("user-123"))
	app.Put("/views/:id", handler.UpdateView)

	tests := []struct {
		name       string
		viewID     string
		payload    interface{}
		setupMocks func()
		wantStatus int
	}{
		{
			name:       "missing viewID param",
			viewID:     "",
			payload:    models.UpdateViewRequest{},
			setupMocks: func() {},
			wantStatus: fiber.StatusNotFound,
		},
		{
			name:       "invalid request body",
			viewID:     "v1",
			payload:    "not-json",
			setupMocks: func() {},
			wantStatus: fiber.StatusBadRequest,
		},
		{
			name:    "view not found",
			viewID:  "v2",
			payload: models.UpdateViewRequest{},
			setupMocks: func() {
				mockRepo.On("GetTeamIDByViewID", mock.Anything, "v2").
					Return("team-a", nil)
				mockRepo.On("GetViewByID", mock.Anything, "v2").
					Return([]db.View{}, nil)
			},
			wantStatus: fiber.StatusNotFound,
		},
		{
			name:    "update name only",
			viewID:  "v3",
			payload: models.UpdateViewRequest{Name: "Updated View"},
			setupMocks: func() {
				sqlMock.ExpectBegin() 

				mockRepo.On("GetTeamIDByViewID", mock.Anything, "v3").Return("team-a", nil)
				mockRepo.On("GetViewByID", mock.Anything, "v3").Return([]db.View{{ID: "v3", Name: "Old"}}, nil)
				mockRepo.On("UpdateViewName", mock.Anything, "v3", "Updated View").Return(nil)
			},
			wantStatus: fiber.StatusOK,
		},
		{
			name:       "missing viewID",
			viewID:     "undefined",
			payload:    models.UpdateViewRequest{},
			setupMocks: func() {},
			wantStatus: fiber.StatusBadRequest,
		},
		{
			name:    "GetTeamIDByViewID error",
			viewID:  "v101",
			payload: models.UpdateViewRequest{},
			setupMocks: func() {
				mockRepo.On("GetTeamIDByViewID", mock.Anything, "v101").
					Return("", errors.New("teamID error"))
			},
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name:    "GetViewByID error",
			viewID:  "v102",
			payload: models.UpdateViewRequest{},
			setupMocks: func() {
				mockRepo.On("GetTeamIDByViewID", mock.Anything, "v102").
					Return("team-X", nil)
				mockRepo.On("GetViewByID", mock.Anything, "v102").
					Return(nil, errors.New("view fetch fail"))
			},
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name:    "tx rollback triggered by UpdateViewName fail",
			viewID:  "v999",
			payload: models.UpdateViewRequest{Name: "Broken"},
			setupMocks: func() {
				mockRepo.On("GetTeamIDByViewID", mock.Anything, "v999").
					Return("team-X", nil)

				mockRepo.On("GetViewByID", mock.Anything, "v999").
					Return([]db.View{{ID: "v999"}}, nil)

				sqlMock.ExpectBegin()    
				sqlMock.ExpectRollback() 

				mockRepo.On("UpdateViewName", mock.Anything, "v999", "Broken").
					Return(errors.New("fail to update"))
			},
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name:    "update team_id fail",
			viewID:  "v-team-fail",
			payload: models.UpdateViewRequest{TeamID: "team-x"},
			setupMocks: func() {

				mockRepo.On("GetTeamIDByViewID", mock.Anything, "v-team-fail").
					Return("team-a", nil)

				mockRepo.On("GetViewByID", mock.Anything, "v-team-fail").
					Return([]db.View{{ID: "v-team-fail"}}, nil)

				sqlMock.ExpectBegin()    
				sqlMock.ExpectRollback() 

				mockRepo.On("UpdateViewTeamID", mock.Anything, "v-team-fail", "team-x").
					Return(errors.New("update team_id failed"))
			},
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name:    "remove group_by fail",
			viewID:  "v201",
			payload: models.UpdateViewRequest{GroupBys: []string{"status"}},
			setupMocks: func() {
				mockRepo.On("GetTeamIDByViewID", mock.Anything, "v201").Return("team-X", nil)
				mockRepo.On("GetViewByID", mock.Anything, "v201").Return([]db.View{{ID: "v201"}}, nil)

				sqlMock.ExpectBegin()
				sqlMock.ExpectRollback()

				mockRepo.On("RemoveGroupByFromView", mock.Anything, "v201").
					Return(errors.New("groupBy remove fail"))
			},
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name:    "remove issues fail",
			viewID:  "v202",
			payload: models.UpdateViewRequest{GroupBys: []string{"status"}},
			setupMocks: func() {
				mockRepo.On("GetTeamIDByViewID", mock.Anything, "v202").Return("team-X", nil)
				mockRepo.On("GetViewByID", mock.Anything, "v202").Return([]db.View{{ID: "v202"}}, nil)

				sqlMock.ExpectBegin()
				sqlMock.ExpectRollback()

				mockRepo.On("RemoveGroupByFromView", mock.Anything, "v202").Return(nil)
				mockRepo.On("RemoveIssueFromView", mock.Anything, "v202").
					Return(errors.New("remove issue fail"))
			},
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name:    "add group_by fail",
			viewID:  "v203",
			payload: models.UpdateViewRequest{GroupBys: []string{"status"}},
			setupMocks: func() {
				mockRepo.On("GetTeamIDByViewID", mock.Anything, "v203").Return("team-X", nil)
				mockRepo.On("GetViewByID", mock.Anything, "v203").Return([]db.View{{ID: "v203"}}, nil)

				sqlMock.ExpectBegin()
				sqlMock.ExpectRollback()

				mockRepo.On("RemoveGroupByFromView", mock.Anything, "v203").Return(nil)
				mockRepo.On("RemoveIssueFromView", mock.Anything, "v203").Return(nil)
				mockRepo.On("AddGroupByToView", mock.Anything, mock.Anything).
					Return(errors.New("groupBy insert fail"))
			},
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name:    "query fail",
			viewID:  "v204",
			payload: models.UpdateViewRequest{GroupBys: []string{"status"}},
			setupMocks: func() {
				mockRepo.On("GetTeamIDByViewID", mock.Anything, "v204").Return("team-X", nil)
				mockRepo.On("GetViewByID", mock.Anything, "v204").Return([]db.View{{ID: "v204"}}, nil)

				sqlMock.ExpectBegin()
				sqlMock.ExpectRollback()

				mockRepo.On("RemoveGroupByFromView", mock.Anything, "v204").Return(nil)
				mockRepo.On("RemoveIssueFromView", mock.Anything, "v204").Return(nil)
				mockRepo.On("AddGroupByToView", mock.Anything, mock.Anything).Return(nil)

				sqlMock.ExpectQuery("SELECT .* FROM issues .*").WithArgs("team-X").
					WillReturnError(errors.New("query error"))
			},
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name:    "rows.Err non-nil",
			viewID:  "v205",
			payload: models.UpdateViewRequest{GroupBys: []string{"status"}},
			setupMocks: func() {
				mockRepo.On("GetTeamIDByViewID", mock.Anything, "v205").Return("team-X", nil)
				mockRepo.On("GetViewByID", mock.Anything, "v205").Return([]db.View{{ID: "v205"}}, nil)

				sqlMock.ExpectBegin()
				sqlMock.ExpectRollback()

				mockRepo.On("RemoveGroupByFromView", mock.Anything, "v205").Return(nil)
				mockRepo.On("RemoveIssueFromView", mock.Anything, "v205").Return(nil)
				mockRepo.On("AddGroupByToView", mock.Anything, mock.Anything).Return(nil)

				rows := sqlmock.NewRows([]string{"issue_id"}).AddRow("i1")
				rows.RowError(0, errors.New("row err"))
				sqlMock.ExpectQuery("SELECT .* FROM issues .*").WithArgs("team-X").
					WillReturnRows(rows)
			},
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name:    "insert issue fail",
			viewID:  "v206",
			payload: models.UpdateViewRequest{GroupBys: []string{"status"}},
			setupMocks: func() {
				mockRepo.On("GetTeamIDByViewID", mock.Anything, "v206").Return("team-X", nil)
				mockRepo.On("GetViewByID", mock.Anything, "v206").Return([]db.View{{ID: "v206"}}, nil)

				sqlMock.ExpectBegin()
				sqlMock.ExpectRollback()

				mockRepo.On("RemoveGroupByFromView", mock.Anything, "v206").Return(nil)
				mockRepo.On("RemoveIssueFromView", mock.Anything, "v206").Return(nil)
				mockRepo.On("AddGroupByToView", mock.Anything, mock.Anything).Return(nil)

				sqlMock.ExpectQuery("SELECT .* FROM issues .*").WithArgs("team-X").
					WillReturnRows(sqlmock.NewRows([]string{"issue_id"}).AddRow("i1"))

				mockRepo.On("AddIssueToViewTx", mock.Anything, mock.Anything, mock.Anything).
					Return(errors.New("insert issue fail"))
			},
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name:    "BeginTx error",
			viewID:  "v103",
			payload: models.UpdateViewRequest{},
			setupMocks: func() {
				mockRepo.On("GetTeamIDByViewID", mock.Anything, "v103").
					Return("team-Y", nil)
				mockRepo.On("GetViewByID", mock.Anything, "v103").
					Return([]db.View{{ID: "v103"}}, nil)
				mockDB.Close() 
			},
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name:    "UpdateViewName error",
			viewID:  "v104",
			payload: models.UpdateViewRequest{Name: "Broken Name"},
			setupMocks: func() {
				mockRepo.On("GetTeamIDByViewID", mock.Anything, "v104").
					Return("team-Z", nil)
				mockRepo.On("GetViewByID", mock.Anything, "v104").
					Return([]db.View{{ID: "v104"}}, nil)
				sqlMock.ExpectBegin()
				mockRepo.On("UpdateViewName", mock.Anything, "v104", "Broken Name").
					Return(errors.New("update fail"))
			},
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name:    "QueryContext error",
			viewID:  "v402",
			payload: models.UpdateViewRequest{GroupBys: []string{"status"}},
			setupMocks: func() {
				mockRepo.On("GetTeamIDByViewID", mock.Anything, "v402").
					Return("team-X", nil)
				mockRepo.On("GetViewByID", mock.Anything, "v402").
					Return([]db.View{{ID: "v402"}}, nil)

				sqlMock.ExpectBegin()
				sqlMock.ExpectRollback()

				mockRepo.On("RemoveGroupByFromView", mock.Anything, "v402").Return(nil)
				mockRepo.On("RemoveIssueFromView", mock.Anything, "v402").Return(nil)
				mockRepo.On("AddGroupByToView", mock.Anything, mock.Anything).Return(nil)

				sqlMock.ExpectQuery("SELECT .* FROM issues .*").
					WithArgs("team-X").
					WillReturnError(errors.New("query fail"))
			},
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name:    "rows.Err error",
			viewID:  "v403",
			payload: models.UpdateViewRequest{GroupBys: []string{"status"}},
			setupMocks: func() {
				mockRepo.On("GetTeamIDByViewID", mock.Anything, "v403").
					Return("team-X", nil)
				mockRepo.On("GetViewByID", mock.Anything, "v403").
					Return([]db.View{{ID: "v403"}}, nil)

				sqlMock.ExpectBegin()
				sqlMock.ExpectRollback()

				mockRepo.On("RemoveGroupByFromView", mock.Anything, "v403").Return(nil)
				mockRepo.On("RemoveIssueFromView", mock.Anything, "v403").Return(nil)
				mockRepo.On("AddGroupByToView", mock.Anything, mock.Anything).Return(nil)

				rows := sqlmock.NewRows([]string{"issue_id"}).AddRow("i1")
				rows.RowError(0, errors.New("row scan error"))
				sqlMock.ExpectQuery("SELECT .* FROM issues .*").
					WithArgs("team-X").WillReturnRows(rows)
			},
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name:    "AddIssueToViewTx error",
			viewID:  "v404",
			payload: models.UpdateViewRequest{GroupBys: []string{"status"}},
			setupMocks: func() {
				mockRepo.On("GetTeamIDByViewID", mock.Anything, "v404").
					Return("team-X", nil)
				mockRepo.On("GetViewByID", mock.Anything, "v404").
					Return([]db.View{{ID: "v404"}}, nil)

				sqlMock.ExpectBegin()
				sqlMock.ExpectRollback()

				mockRepo.On("RemoveGroupByFromView", mock.Anything, "v404").Return(nil)
				mockRepo.On("RemoveIssueFromView", mock.Anything, "v404").Return(nil)
				mockRepo.On("AddGroupByToView", mock.Anything, mock.Anything).Return(nil)

				sqlMock.ExpectQuery("SELECT .* FROM issues .*").
					WithArgs("team-X").
					WillReturnRows(sqlmock.NewRows([]string{"issue_id"}).
						AddRow("i1"))

				mockRepo.On("AddIssueToViewTx", mock.Anything, mock.Anything, mock.MatchedBy(
					func(p db.AddIssueToViewParams) bool {
						return p.ViewID == "v404" && p.IssueID == "i1"
					})).Return(errors.New("insert error"))
			},
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name:    "tx.QueryContext error",
			viewID:  "v502",
			payload: models.UpdateViewRequest{GroupBys: []string{"status"}},
			setupMocks: func() {
				mockRepo.On("GetTeamIDByViewID", mock.Anything, "v502").Return("team-x", nil)
				mockRepo.On("GetViewByID", mock.Anything, "v502").Return([]db.View{{ID: "v502"}}, nil)
				sqlMock.ExpectBegin()
				sqlMock.ExpectRollback()
				mockRepo.On("RemoveGroupByFromView", mock.Anything, "v502").Return(nil)
				mockRepo.On("RemoveIssueFromView", mock.Anything, "v502").Return(nil)
				mockRepo.On("AddGroupByToView", mock.Anything, mock.Anything).Return(nil)

				sqlMock.ExpectQuery("SELECT .* FROM issues .*").
					WithArgs("team-x").
					WillReturnError(errors.New("db query fail"))
			},
			wantStatus: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			var body []byte
			if str, ok := tt.payload.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.payload)
			}

			url := "/views/" + tt.viewID
			req := httptest.NewRequest(http.MethodPut, url, bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}

	t.Run("group_by query, scan, and insert issue success", func(t *testing.T) {
		mockRepo := new(mocks.MockViewRepo)
		mockDB, sqlMock, err := sqlmock.New()
		require.NoError(t, err)
		defer mockDB.Close()

		handler := &routes.ViewHandler{
			DB:   mockDB,
			Repo: mockRepo,
		}

		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Put("/views/:id", handler.UpdateView)

		query, err := routes.BuildGroupByQuery("issues", []string{"status"})
		require.NoError(t, err)

		sqlMock.ExpectBegin()
		sqlMock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs("team-x").
			WillReturnRows(sqlmock.NewRows([]string{"issue_id"}).
				AddRow("i1").AddRow("i2"))

		mockRepo.On("GetTeamIDByViewID", mock.Anything, "v700").Return("team-x", nil)
		mockRepo.On("GetViewByID", mock.Anything, "v700").Return([]db.View{{ID: "v700"}}, nil)

		sqlMock.ExpectBegin()
		sqlMock.ExpectCommit()

		mockRepo.On("RemoveGroupByFromView", mock.Anything, "v700").Return(nil)
		mockRepo.On("RemoveIssueFromView", mock.Anything, "v700").Return(nil)
		mockRepo.On("AddGroupByToView", mock.Anything, mock.Anything).Return(nil)

		mockRepo.On("AddIssueToViewTx", mock.Anything, mock.Anything,
			mock.MatchedBy(func(p db.AddIssueToViewParams) bool {
				return p.ViewID == "v700" && (p.IssueID == "i1" || p.IssueID == "i2")
			})).Return(nil).Times(2)

		payload := models.UpdateViewRequest{GroupBys: []string{"status"}}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPut, "/views/v700", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("query grouped issues fail", func(t *testing.T) {
		mockRepo := new(mocks.MockViewRepo)
		mockDB, sqlMock, _ := sqlmock.New()
		defer mockDB.Close()

		handler := &routes.ViewHandler{DB: mockDB, Repo: mockRepo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Put("/views/:id", handler.UpdateView)

		query, _ := routes.BuildGroupByQuery("issues", []string{"status"})

		mockRepo.On("GetTeamIDByViewID", mock.Anything, "v800").Return("team-x", nil)
		mockRepo.On("GetViewByID", mock.Anything, "v800").Return([]db.View{{ID: "v800"}}, nil)

		sqlMock.ExpectBegin()
		sqlMock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs("team-x").
			WillReturnError(fmt.Errorf("db query error"))
		sqlMock.ExpectRollback()

		mockRepo.On("RemoveGroupByFromView", mock.Anything, "v800").Return(nil)
		mockRepo.On("RemoveIssueFromView", mock.Anything, "v800").Return(nil)
		mockRepo.On("AddGroupByToView", mock.Anything, mock.Anything).Return(nil)

		payload := models.UpdateViewRequest{GroupBys: []string{"status"}}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPut, "/views/v800", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("row scan error, but continues", func(t *testing.T) {
		mockRepo := new(mocks.MockViewRepo)
		mockDB, sqlMock, _ := sqlmock.New()
		defer mockDB.Close()

		handler := &routes.ViewHandler{DB: mockDB, Repo: mockRepo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Put("/views/:id", handler.UpdateView)

		query, _ := routes.BuildGroupByQuery("issues", []string{"status"})

		mockRepo.On("GetTeamIDByViewID", mock.Anything, "v801").Return("team-x", nil)
		mockRepo.On("GetViewByID", mock.Anything, "v801").Return([]db.View{{ID: "v801"}}, nil)

		sqlMock.ExpectBegin()
		sqlMock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs("team-x").
			WillReturnRows(sqlmock.NewRows([]string{"issue_id"}).AddRow(nil)) 
		sqlMock.ExpectCommit()

		mockRepo.On("RemoveGroupByFromView", mock.Anything, "v801").Return(nil)
		mockRepo.On("RemoveIssueFromView", mock.Anything, "v801").Return(nil)
		mockRepo.On("AddGroupByToView", mock.Anything, mock.Anything).Return(nil)

		payload := models.UpdateViewRequest{GroupBys: []string{"status"}}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPut, "/views/v801", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode) 
		mockRepo.AssertExpectations(t)
	})

	t.Run("rows.Err after scan", func(t *testing.T) {
		mockRepo := new(mocks.MockViewRepo)
		mockDB, sqlMock, _ := sqlmock.New()
		defer mockDB.Close()

		handler := &routes.ViewHandler{DB: mockDB, Repo: mockRepo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Put("/views/:id", handler.UpdateView)

		query, _ := routes.BuildGroupByQuery("issues", []string{"status"})

		rows := sqlmock.NewRows([]string{"issue_id"}).
			AddRow("i1").
			RowError(0, fmt.Errorf("row error"))

		mockRepo.On("GetTeamIDByViewID", mock.Anything, "v802").Return("team-x", nil)
		mockRepo.On("GetViewByID", mock.Anything, "v802").Return([]db.View{{ID: "v802"}}, nil)

		sqlMock.ExpectBegin()
		sqlMock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs("team-x").
			WillReturnRows(rows)
		sqlMock.ExpectRollback()

		mockRepo.On("RemoveGroupByFromView", mock.Anything, "v802").Return(nil)
		mockRepo.On("RemoveIssueFromView", mock.Anything, "v802").Return(nil)
		mockRepo.On("AddGroupByToView", mock.Anything, mock.Anything).Return(nil)

		payload := models.UpdateViewRequest{GroupBys: []string{"status"}}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPut, "/views/v802", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("rows.Err after scan", func(t *testing.T) {
		mockRepo := new(mocks.MockViewRepo)
		mockDB, sqlMock, _ := sqlmock.New()
		defer mockDB.Close()

		handler := &routes.ViewHandler{DB: mockDB, Repo: mockRepo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Put("/views/:id", handler.UpdateView)

		query, _ := routes.BuildGroupByQuery("issues", []string{"status"})

		rows := sqlmock.NewRows([]string{"issue_id"}).
			AddRow("i1").
			RowError(0, fmt.Errorf("row error"))

		mockRepo.On("GetTeamIDByViewID", mock.Anything, "v802").Return("team-x", nil)
		mockRepo.On("GetViewByID", mock.Anything, "v802").Return([]db.View{{ID: "v802"}}, nil)

		sqlMock.ExpectBegin()
		sqlMock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs("team-x").
			WillReturnRows(rows)
		sqlMock.ExpectRollback()

		mockRepo.On("RemoveGroupByFromView", mock.Anything, "v802").Return(nil)
		mockRepo.On("RemoveIssueFromView", mock.Anything, "v802").Return(nil)
		mockRepo.On("AddGroupByToView", mock.Anything, mock.Anything).Return(nil)

		payload := models.UpdateViewRequest{GroupBys: []string{"status"}}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPut, "/views/v802", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("AddIssueToViewTx fail", func(t *testing.T) {
		mockRepo := new(mocks.MockViewRepo)
		mockDB, sqlMock, _ := sqlmock.New()
		defer mockDB.Close()

		handler := &routes.ViewHandler{DB: mockDB, Repo: mockRepo}
		app := fiber.New()
		app.Use(withUserID("user-123"))
		app.Put("/views/:id", handler.UpdateView)

		query, _ := routes.BuildGroupByQuery("issues", []string{"status"})

		mockRepo.On("GetTeamIDByViewID", mock.Anything, "v803").Return("team-x", nil)
		mockRepo.On("GetViewByID", mock.Anything, "v803").Return([]db.View{{ID: "v803"}}, nil)

		sqlMock.ExpectBegin()
		sqlMock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs("team-x").
			WillReturnRows(sqlmock.NewRows([]string{"issue_id"}).AddRow("i123"))
		sqlMock.ExpectRollback()

		mockRepo.On("RemoveGroupByFromView", mock.Anything, "v803").Return(nil)
		mockRepo.On("RemoveIssueFromView", mock.Anything, "v803").Return(nil)
		mockRepo.On("AddGroupByToView", mock.Anything, mock.Anything).Return(nil)

		mockRepo.On("AddIssueToViewTx", mock.Anything, mock.Anything,
			mock.MatchedBy(func(p db.AddIssueToViewParams) bool {
				return p.ViewID == "v803" && p.IssueID == "i123"
			})).Return(fmt.Errorf("insert error"))

		payload := models.UpdateViewRequest{GroupBys: []string{"status"}}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPut, "/views/v803", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		mockRepo.AssertExpectations(t)
	})

}
