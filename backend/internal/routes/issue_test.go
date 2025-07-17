package routes_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
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

func ptr[T any](v T) *T { return &v }
func mustJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

func TestNewIssueHandler(t *testing.T) {
	db := &sql.DB{}
	mockRepo := new(mocks.MockIssueRepo)
	projectRepo := new(mocks.MockProjectRepo)
	teamRepo := new(mocks.MockTeamRepository)

	handler := routes.NewIssueHandler(db, mockRepo, teamRepo, projectRepo)

	assert.Equal(t, db, handler.DB)
	assert.Equal(t, mockRepo, handler.Repo)
	assert.Equal(t, teamRepo, handler.TeamRepo)
	assert.Equal(t, projectRepo, handler.ProjectRepo)
}

func TestCreateIssue(t *testing.T) {
	app := fiber.New()
	db := &sql.DB{}
	mockRepo := new(mocks.MockIssueRepo)
	mockTeamRepo := new(mocks.MockTeamRepository)
	mockProjRepo := new(mocks.MockProjectRepo)
	handler := routes.IssueHandler{db, mockRepo, mockTeamRepo, mockProjRepo}

	app.Use(withUserID("user-123"))
	app.Post("/issues", handler.CreateIssue)

	type mocks struct {
		repo    func()
		team    func()
		project func()
	}

	tests := []struct {
		name       string
		body       *models.IssueCreate
		rawBody    []byte
		setupMocks mocks
		wantStatus int
	}{
		{
			name:       "Create Issue Successfully",
			body:       &models.IssueCreate{TeamID: "team-1", Title: "Test Issue"},
			wantStatus: fiber.StatusCreated,
			setupMocks: mocks{
				repo: func() {
					mockRepo.On("CreateIssue", mock.Anything, mock.AnythingOfType("db.CreateIssueParams")).Return(nil)
				},
				team: func() {
					mockTeamRepo.On("IsTeamExists", mock.Anything, "team-1").Return(true, nil)
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-1", "user-123").Return(true, nil)
				},
				project: func() {},
			},
		},
		{
			name:       "Create Issue Failed (repo error)",
			body:       &models.IssueCreate{TeamID: "team-1", Title: "Test Issue"},
			wantStatus: fiber.StatusInternalServerError,
			setupMocks: mocks{
				repo: func() {
					mockRepo.On("CreateIssue", mock.Anything, mock.AnythingOfType("db.CreateIssueParams")).Return(assert.AnError)
				},
				team: func() {
					mockTeamRepo.On("IsTeamExists", mock.Anything, "team-1").Return(true, nil)
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-1", "user-123").Return(true, nil)
				},
				project: func() {},
			},
		},
		{
			name:       "invalid JSON body",
			rawBody:    []byte(`{invalid`),
			body:       nil,
			wantStatus: fiber.StatusBadRequest,
			setupMocks: mocks{repo: func() {}, team: func() {}, project: func() {}},
		},
		{
			name:       "user not in team",
			body:       &models.IssueCreate{TeamID: "team-1", Title: "Unauthorized"},
			wantStatus: fiber.StatusForbidden,
			setupMocks: mocks{
				repo: func() {},
				team: func() {
					mockTeamRepo.On("IsTeamExists", mock.Anything, "team-1").Return(true, nil)
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-1", "user-123").Return(false, nil)
				},
				project: func() {},
			},
		},
		{
			name:       "team not found",
			body:       &models.IssueCreate{TeamID: "invalid-team", Title: "No Team"},
			wantStatus: fiber.StatusBadRequest,
			setupMocks: mocks{
				repo: func() {},
				team: func() {
					mockTeamRepo.On("IsTeamExists", mock.Anything, "invalid-team").Return(false, nil)
				},
				project: func() {},
			},
		},
		{
			name:       "error on checking team member",
			body:       &models.IssueCreate{TeamID: "team-1", Title: "Error Member"},
			wantStatus: fiber.StatusInternalServerError,
			setupMocks: mocks{
				repo: func() {},
				team: func() {
					mockTeamRepo.On("IsTeamExists", mock.Anything, "team-1").Return(true, nil)
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-1", "user-123").Return(false, errors.New("mock error"))
				},
				project: func() {},
			},
		},
		{
			name:       "project not found",
			body:       &models.IssueCreate{TeamID: "team-1", Title: "Missing Project", ProjectID: ptr("proj-x")},
			wantStatus: fiber.StatusBadRequest,
			setupMocks: mocks{
				repo: func() {},
				team: func() {
					mockTeamRepo.On("IsTeamExists", mock.Anything, "team-1").Return(true, nil)
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-1", "user-123").Return(true, nil)
				},
				project: func() {
					mockProjRepo.On("IsProjectExists", mock.Anything, "proj-x").Return(false, nil)
				},
			},
		},
		{
			name:       "assignee not in team",
			body:       &models.IssueCreate{TeamID: "team-1", Title: "Assignee Fail", Assignee: &[]string{"user-x"}},
			wantStatus: fiber.StatusCreated,
			setupMocks: mocks{
				repo: func() {
					mockRepo.On("CreateIssue", mock.Anything, mock.AnythingOfType("db.CreateIssueParams")).Return(nil)
				},
				team: func() {
					mockTeamRepo.On("IsTeamExists", mock.Anything, "team-1").Return(true, nil)
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-1", "user-123").Return(true, nil)
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-1", "user-x").Return(false, nil)
				},
				project: func() {},
			},
		},
		{
			name:       "error adding assignee",
			body:       &models.IssueCreate{TeamID: "team-1", Title: "Assignee DB Error", Assignee: &[]string{"user-y"}},
			wantStatus: fiber.StatusCreated,
			setupMocks: mocks{
				repo: func() {
					mockRepo.On("CreateIssue", mock.Anything, mock.AnythingOfType("db.CreateIssueParams")).Return(nil)
					mockRepo.On("AddAssigneeToIssue", mock.Anything, mock.Anything).Return(assert.AnError)
				},
				team: func() {
					mockTeamRepo.On("IsTeamExists", mock.Anything, "team-1").Return(true, nil)
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-1", "user-123").Return(true, nil)
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-1", "user-y").Return(true, nil)
				},
				project: func() {},
			},
		},
		{
			name:       "validation failed",
			body:       &models.IssueCreate{TeamID: "", Title: ""},
			wantStatus: fiber.StatusBadRequest,
			setupMocks: mocks{
				repo:    func() {},
				team:    func() {},
				project: func() {},
			},
		},
		{
			name:       "error checking assignee membership",
			body:       &models.IssueCreate{TeamID: "team-1", Title: "With Assignee Error", Assignee: &[]string{"user-z"}},
			wantStatus: fiber.StatusCreated,
			setupMocks: mocks{
				repo: func() {
					mockRepo.On("CreateIssue", mock.Anything, mock.AnythingOfType("db.CreateIssueParams")).Return(nil)
				},
				team: func() {
					mockTeamRepo.On("IsTeamExists", mock.Anything, "team-1").Return(true, nil)
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-1", "user-123").Return(true, nil)
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-1", "user-z").Return(false, errors.New("mock error"))
				},
				project: func() {},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.body != nil {
				tt.rawBody = mustJSON(tt.body)
			}

			tt.setupMocks.repo()
			tt.setupMocks.team()
			tt.setupMocks.project()

			req := httptest.NewRequest(http.MethodPost, "/issues", bytes.NewReader(tt.rawBody))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			mockRepo.ExpectedCalls = nil
			mockTeamRepo.ExpectedCalls = nil
			mockProjRepo.ExpectedCalls = nil
		})
	}
}

func TestUpdateIssue(t *testing.T) {
	app := fiber.New()
	mockDB, sqlMock, err := sqlmock.New()
	require.NoError(t, err)

	mockRepo := new(mocks.MockIssueRepo)
	mockTeamRepo := new(mocks.MockTeamRepository)
	mockProjRepo := new(mocks.MockProjectRepo)

	handler := routes.IssueHandler{
		DB:          mockDB,
		Repo:        mockRepo,
		TeamRepo:    mockTeamRepo,
		ProjectRepo: mockProjRepo,
	}

	app.Use(withUserID("user-123"))
	app.Put("/issues/:id", handler.UpdateIssue)

	type mocks struct {
		repo  func()
		team  func()
		query func()
	}

	tests := []struct {
		name       string
		req        *models.UpdateIssueRequest
		rawBody    []byte
		issueID    string
		setupMocks mocks
		wantStatus int
	}{
		{
			name:       "success - authorized owner with update query",
			issueID:    "issue-1",
			req:        &models.UpdateIssueRequest{Title: ptr("Updated Title")},
			wantStatus: fiber.StatusOK,
			setupMocks: mocks{
				repo: func() {
					mockRepo.On("GetIssueByID", mock.Anything, "issue-1").
						Return(db.Issue{ID: "issue-1", TeamID: "team-1", OwnerID: "user-123"}, nil)
				},
				team: func() {
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-1", "user-123").
						Return(true, nil)
				},
				query: func() {
					sqlMock.ExpectExec("UPDATE issues SET title = .*").
						WithArgs("Updated Title", "issue-1").
						WillReturnResult(sqlmock.NewResult(1, 1))
				},
			},
		},
		{
			name:       "update query failed",
			issueID:    "issue-err",
			req:        &models.UpdateIssueRequest{Title: ptr("Failed")},
			wantStatus: fiber.StatusInternalServerError,
			setupMocks: mocks{
				repo: func() {
					mockRepo.On("GetIssueByID", mock.Anything, "issue-err").
						Return(db.Issue{ID: "issue-err", TeamID: "team-X", OwnerID: "user-123"}, nil)
				},
				team: func() {
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-X", "user-123").
						Return(true, nil)
				},
				query: func() {
					sqlMock.ExpectExec("UPDATE issues SET title = .*").
						WithArgs("Failed", "issue-err").
						WillReturnError(errors.New("query failed"))
				},
			},
		},
		{
			name:       " invalid JSON",
			issueID:    "issue-1",
			rawBody:    []byte(`{bad}`),
			wantStatus: fiber.StatusBadRequest,
			setupMocks: mocks{repo: func() {}, team: func() {}, query: func() {}},
		},
		{
			name:       "issue ID is 'undefined'",
			issueID:    "undefined",
			req:        &models.UpdateIssueRequest{Title: ptr("Title")},
			wantStatus: fiber.StatusBadRequest,
			setupMocks: mocks{repo: func() {}, team: func() {}, query: func() {}},
		},
		{
			name:       "issue not found",
			issueID:    "issue-x",
			req:        &models.UpdateIssueRequest{Title: ptr("Title")},
			wantStatus: fiber.StatusNotFound,
			setupMocks: mocks{
				repo: func() {
					mockRepo.On("GetIssueByID", mock.Anything, "issue-x").
						Return(db.Issue{}, sql.ErrNoRows)
				},
				team: func() {}, query: func() {},
			},
		},
		{
			name:       "repo error",
			issueID:    "issue-2",
			req:        &models.UpdateIssueRequest{Title: ptr("Oops")},
			wantStatus: fiber.StatusInternalServerError,
			setupMocks: mocks{
				repo: func() {
					mockRepo.On("GetIssueByID", mock.Anything, "issue-2").
						Return(db.Issue{}, errors.New("unexpected error"))
				},
				team: func() {}, query: func() {},
			},
		},
		{
			name:       "unauthorized",
			issueID:    "issue-3",
			req:        &models.UpdateIssueRequest{Title: ptr("Denied")},
			wantStatus: fiber.StatusForbidden,
			setupMocks: mocks{
				repo: func() {
					mockRepo.On("GetIssueByID", mock.Anything, "issue-3").
						Return(db.Issue{ID: "issue-3", TeamID: "team-1", OwnerID: "someone-else"}, nil)
				},
				team: func() {
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-1", "user-123").
						Return(false, nil)
				},
				query: func() {},
			},
		},
		{
			name:       "error checking team membership",
			issueID:    "issue-4",
			req:        &models.UpdateIssueRequest{Title: ptr("TeamErr")},
			wantStatus: fiber.StatusInternalServerError,
			setupMocks: mocks{
				repo: func() {
					mockRepo.On("GetIssueByID", mock.Anything, "issue-4").
						Return(db.Issue{ID: "issue-4", TeamID: "team-err", OwnerID: "owner-1"}, nil)
				},
				team: func() {
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-err", "user-123").
						Return(false, errors.New("mock error"))
				},
				query: func() {},
			},
		},
		{
			name:       "AddAssignee: success",
			issueID:    "issue-add-ok",
			req:        &models.UpdateIssueRequest{AddAssignee: &[]string{"user-x"}},
			wantStatus: fiber.StatusOK,
			setupMocks: mocks{
				repo: func() {
					mockRepo.On("GetIssueByID", mock.Anything, "issue-add-ok").
						Return(db.Issue{ID: "issue-add-ok", TeamID: "team-1", OwnerID: "user-123"}, nil)
					mockRepo.On("AddAssigneeToIssue", mock.Anything, db.AddAssigneeToIssueParams{
						IssueID: "issue-add-ok",
						UserID:  "user-x",
					}).Return(nil)
				},
				team: func() {
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-1", "user-123").
						Return(true, nil)
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-1", "user-x").
						Return(true, nil)
				},
				query: func() {},
			},
		},
		{
			name:       "AddAssignee: error checking membership",
			issueID:    "issue-add-check-error",
			req:        &models.UpdateIssueRequest{AddAssignee: &[]string{"user-x"}},
			wantStatus: fiber.StatusOK,
			setupMocks: mocks{
				repo: func() {
					mockRepo.On("GetIssueByID", mock.Anything, "issue-add-check-error").
						Return(db.Issue{ID: "issue-add-check-error", TeamID: "team-a", OwnerID: "user-123"}, nil)
				},
				team: func() {
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-a", "user-123").
						Return(true, nil)
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-a", "user-x").
						Return(false, errors.New("check error"))
				},
				query: func() {},
			},
		},
		{
			name:       "RemoveAssignee: error removing",
			issueID:    "remove-error",
			req:        &models.UpdateIssueRequest{RemoveAssignee: &[]string{"user-x"}},
			wantStatus: fiber.StatusOK,
			setupMocks: mocks{
				repo: func() {
					mockRepo.On("GetIssueByID", mock.Anything, "remove-error").
						Return(db.Issue{ID: "remove-error", TeamID: "team-b", OwnerID: "user-123"}, nil)
					mockRepo.On("RemoveAssigneeFromIssue", mock.Anything, mock.Anything).
						Return(assert.AnError)
				},
				team: func() {
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-b", "user-123").
						Return(true, nil)
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-b", "user-x").
						Return(true, nil)
				},
				query: func() {},
			},
		},
		{
			name:       "AddAssignee: not member of team",
			issueID:    "issue-add-notmember",
			req:        &models.UpdateIssueRequest{AddAssignee: &[]string{"user-x"}},
			wantStatus: fiber.StatusOK,
			setupMocks: mocks{
				repo: func() {
					mockRepo.On("GetIssueByID", mock.Anything, "issue-add-notmember").
						Return(db.Issue{ID: "issue-add-notmember", TeamID: "team-z", OwnerID: "user-123"}, nil)
				},
				team: func() {
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-z", "user-123").
						Return(true, nil)
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-z", "user-x").
						Return(false, nil)
				},
				query: func() {},
			},
		},
		{
			name:       "AddAssignee: error while adding",
			issueID:    "issue-add-fail",
			req:        &models.UpdateIssueRequest{AddAssignee: &[]string{"user-x"}},
			wantStatus: fiber.StatusOK,
			setupMocks: mocks{
				repo: func() {
					mockRepo.On("GetIssueByID", mock.Anything, "issue-add-fail").
						Return(db.Issue{ID: "issue-add-fail", TeamID: "team-z", OwnerID: "user-123"}, nil)
					mockRepo.On("AddAssigneeToIssue", mock.Anything, mock.Anything).
						Return(assert.AnError)
				},
				team: func() {
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-z", "user-123").
						Return(true, nil)
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-z", "user-x").
						Return(true, nil)
				},
				query: func() {},
			},
		},
		{
			name:       "RemoveAssignee: not member of team",
			issueID:    "issue-remove-notmember",
			req:        &models.UpdateIssueRequest{RemoveAssignee: &[]string{"user-y"}},
			wantStatus: fiber.StatusOK,
			setupMocks: mocks{
				repo: func() {
					mockRepo.On("GetIssueByID", mock.Anything, "issue-remove-notmember").
						Return(db.Issue{ID: "issue-remove-notmember", TeamID: "team-y", OwnerID: "user-123"}, nil)
				},
				team: func() {
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-y", "user-123").
						Return(true, nil)
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-y", "user-y").
						Return(false, nil)
				},
				query: func() {},
			},
		},
		{
			name:       "RemoveAssignee: error while removing",
			issueID:    "issue-remove-fail",
			req:        &models.UpdateIssueRequest{RemoveAssignee: &[]string{"user-y"}},
			wantStatus: fiber.StatusOK,
			setupMocks: mocks{
				repo: func() {
					mockRepo.On("GetIssueByID", mock.Anything, "issue-remove-fail").
						Return(db.Issue{ID: "issue-remove-fail", TeamID: "team-y", OwnerID: "user-123"}, nil)
					mockRepo.On("RemoveAssigneeFromIssue", mock.Anything, mock.Anything).
						Return(assert.AnError)
				},
				team: func() {
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-y", "user-123").
						Return(true, nil)
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-y", "user-y").
						Return(true, nil)
				},
				query: func() {},
			},
		},
		{
			name:       "RemoveAssignee: success",
			issueID:    "issue-remove-ok",
			req:        &models.UpdateIssueRequest{RemoveAssignee: &[]string{"user-a"}},
			wantStatus: fiber.StatusOK,
			setupMocks: mocks{
				repo: func() {
					mockRepo.On("GetIssueByID", mock.Anything, "issue-remove-ok").
						Return(db.Issue{ID: "issue-remove-ok", TeamID: "team-x", OwnerID: "user-123"}, nil)
					mockRepo.On("RemoveAssigneeFromIssue", mock.Anything, mock.Anything).
						Return(nil)
				},
				team: func() {
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-x", "user-123").
						Return(true, nil)
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-x", "user-a").
						Return(true, nil)
				},
				query: func() {},
			},
		},
		{
			name:       "RemoveAssignee: error checking membership",
			issueID:    "issue-remove-check-error",
			req:        &models.UpdateIssueRequest{RemoveAssignee: &[]string{"user-z"}},
			wantStatus: fiber.StatusOK,
			setupMocks: mocks{
				repo: func() {
					mockRepo.On("GetIssueByID", mock.Anything, "issue-remove-check-error").
						Return(db.Issue{ID: "issue-remove-check-error", TeamID: "team-z", OwnerID: "user-123"}, nil)
				},
				team: func() {
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-z", "user-123").
						Return(true, nil)
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-z", "user-z").
						Return(false, errors.New("mock error"))
				},
				query: func() {},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req != nil {
				tt.rawBody = mustJSON(tt.req)
			}

			tt.setupMocks.repo()
			tt.setupMocks.team()
			tt.setupMocks.query()

			req := httptest.NewRequest(http.MethodPut, "/issues/"+tt.issueID, bytes.NewReader(tt.rawBody))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)

			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
			require.NoError(t, sqlMock.ExpectationsWereMet())

			mockRepo.ExpectedCalls = nil
			mockTeamRepo.ExpectedCalls = nil
			mockProjRepo.ExpectedCalls = nil
		})
	}
}

func TestDeleteIssue(t *testing.T) {
	app := fiber.New()
	mockRepo := new(mocks.MockIssueRepo)
	mockTeamRepo := new(mocks.MockTeamRepository)
	handler := routes.IssueHandler{
		DB:       nil,
		Repo:     mockRepo,
		TeamRepo: mockTeamRepo,
	}

	app.Use(withUserID("user-123"))
	app.Delete("/issues/:id", handler.DeleteIssue)

	type mocks struct {
		repo func()
		team func()
	}

	tests := []struct {
		name       string
		issueID    string
		setupMocks mocks
		wantStatus int
	}{
		{
			name:    "Success",
			issueID: "success-id",
			setupMocks: mocks{
				repo: func() {
					mockRepo.On("GetIssueByID", mock.Anything, "success-id").
						Return(db.Issue{ID: "success-id", TeamID: "team-1", OwnerID: "user-123"}, nil)
					mockRepo.On("DeleteIssue", mock.Anything, "success-id").
						Return(nil)
				},
				team: func() {
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-1", "user-123").
						Return(true, nil)
				},
			},
			wantStatus: fiber.StatusOK,
		},
		{
			name:       "Missing ID",
			issueID:    "undefined",
			setupMocks: mocks{repo: func() {}, team: func() {}},
			wantStatus: fiber.StatusBadRequest,
		},
		{
			name:    "🔍 Issue not found",
			issueID: "not-found-id",
			setupMocks: mocks{
				repo: func() {
					mockRepo.On("GetIssueByID", mock.Anything, "not-found-id").
						Return(db.Issue{}, sql.ErrNoRows)
				},
				team: func() {},
			},
			wantStatus: fiber.StatusNotFound,
		},
		{
			name:    "DB error getting issue",
			issueID: "db-error-id",
			setupMocks: mocks{
				repo: func() {
					mockRepo.On("GetIssueByID", mock.Anything, "db-error-id").
						Return(db.Issue{}, errors.New("db error"))
				},
				team: func() {},
			},
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name:    "Unauthorized",
			issueID: "unauth-id",
			setupMocks: mocks{
				repo: func() {
					mockRepo.On("GetIssueByID", mock.Anything, "unauth-id").
						Return(db.Issue{ID: "unauth-id", TeamID: "team-2", OwnerID: "someone-else"}, nil)
				},
				team: func() {
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-2", "user-123").
						Return(false, nil)
				},
			},
			wantStatus: fiber.StatusForbidden,
		},
		{
			name:    "Error checking team membership",
			issueID: "team-error-id",
			setupMocks: mocks{
				repo: func() {
					mockRepo.On("GetIssueByID", mock.Anything, "team-error-id").
						Return(db.Issue{ID: "team-error-id", TeamID: "team-err", OwnerID: "user-123"}, nil)
				},
				team: func() {
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-err", "user-123").
						Return(false, errors.New("check error"))
				},
			},
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name:    "Delete error",
			issueID: "delete-error-id",
			setupMocks: mocks{
				repo: func() {
					mockRepo.On("GetIssueByID", mock.Anything, "delete-error-id").
						Return(db.Issue{ID: "delete-error-id", TeamID: "team-x", OwnerID: "user-123"}, nil)
					mockRepo.On("DeleteIssue", mock.Anything, "delete-error-id").
						Return(errors.New("delete error"))
				},
				team: func() {
					mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-x", "user-123").
						Return(true, nil)
				},
			},
			wantStatus: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks.repo()
			tt.setupMocks.team()

			url := "/issues/" + tt.issueID
			req := httptest.NewRequest(http.MethodDelete, url, nil)
			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			mockRepo.ExpectedCalls = nil
			mockTeamRepo.ExpectedCalls = nil
		})
	}
}

func TestGetIssuesByUserID(t *testing.T) {
	app := fiber.New()
	mockRepo := new(mocks.MockIssueRepo)
	handler := routes.IssueHandler{
		DB:          nil,
		Repo:        mockRepo,
		TeamRepo:    nil,
		ProjectRepo: nil,
	}

	app.Use(withUserID("user-123"))
	app.Get("/issues/user", handler.GetIssuesByUserID)

	tests := []struct {
		name       string
		setupMocks func()
		wantStatus int
	}{
		{
			name: "success",
			setupMocks: func() {
				mockRepo.On("GetIssueByUserID", mock.Anything, "user-123").
					Return([]db.Issue{
						{ID: "i1", Title: "Test 1"},
						{ID: "i2", Title: "Test 2"},
					}, nil)
			},
			wantStatus: fiber.StatusOK,
		},
		{
			name: "repo error",
			setupMocks: func() {
				mockRepo.On("GetIssueByUserID", mock.Anything, "user-123").
					Return([]db.Issue{}, errors.New("unexpected error"))
			},
			wantStatus: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil

			tt.setupMocks()

			req := httptest.NewRequest(http.MethodGet, "/issues/user", nil)
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)
		})
	}
}
