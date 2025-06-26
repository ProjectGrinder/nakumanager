package ws_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/nack098/nakumanager/internal/db"
	models "github.com/nack098/nakumanager/internal/models"
	mocks "github.com/nack098/nakumanager/internal/routes/mock_repo"
	"github.com/nack098/nakumanager/internal/ws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateIssueHandler_Success(t *testing.T) {
	mockIssueRepo := new(mocks.MockIssueRepo)
	mockTeamRepo := new(mocks.MockTeamRepo)
	mockUserRepo := new(mocks.MockUserRepo)
	mockConn := new(MockConn)

	handler := ws.NewWSHandler(
		nil,
		mockTeamRepo,
		nil,
		mockIssueRepo,
		mockUserRepo,
		nil,
	)

	userID := "user-123"
	issueID := "issue-1"

	mockConn.On("Locals", "userID").Return(userID)

	existingIssue := db.Issue{
		ID:      issueID,
		OwnerID: userID,
	}
	mockIssueRepo.On("GetIssueByID", mock.Anything, issueID).Return(existingIssue, nil)
	mockTeamRepo.On("IsMemberInTeam", mock.Anything, "team-1", userID).Return(true, nil)
	mockUserRepo.On("GetUserByID", mock.Anything, "assignee-123").Return(db.User{}, nil)
	mockIssueRepo.On("AddAssigneeToIssue", mock.Anything, db.AddAssigneeToIssueParams{
		IssueID: issueID,
		UserID:  "assignee-123",
	}).Return(nil)
	mockIssueRepo.On("UpdateIssue", mock.Anything, mock.AnythingOfType("db.UpdateIssueParams")).Return(nil)

	called := false
	handler.BroadcastFunc = func(data interface{}) {
		broadcastData := data.(map[string]interface{})
		assert.Equal(t, "issue_updated", broadcastData["event"])
		assert.Equal(t, issueID, broadcastData["data"].(models.EditIssue).ID)
		called = true
	}

	title := "Issue Title"
	status := "open"
	teamID := "team-1"
	ownerID := userID
	assignee := "assignee-123"
	priority := "high"
	content := "Important content"
	projectID := "project-456"
	label := "bug"
	start := time.Now()
	end := start.Add(24 * time.Hour)

	issue := models.EditIssue{
		ID:        issueID,
		Title:     &title,
		Status:    &status,
		TeamID:    &teamID,
		OwnerID:   &ownerID,
		Assignee:  &assignee,
		Priority:  &priority,
		Content:   &content,
		ProjectID: &projectID,
		Label:     &label,
		StartDate: &start,
		EndDate:   &end,
	}

	data, _ := json.Marshal(issue)
	handler.UpdateIssueHandler(mockConn, data)

	assert.True(t, called, "Broadcast should be called")
	mockConn.AssertExpectations(t)
	mockIssueRepo.AssertExpectations(t)
	mockTeamRepo.AssertExpectations(t)
	mockUserRepo.AssertExpectations(t)
}
