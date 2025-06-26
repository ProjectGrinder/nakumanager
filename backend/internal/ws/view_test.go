package ws_test

import (
	"database/sql"
	"encoding/json"
	"testing"

	"github.com/nack098/nakumanager/internal/db"
	models "github.com/nack098/nakumanager/internal/models"
	mocks "github.com/nack098/nakumanager/internal/routes/mock_repo"
	"github.com/nack098/nakumanager/internal/ws"
	"github.com/stretchr/testify/mock"
)

func TestUpdateViewHandler_Success(t *testing.T) {
	mockViewRepo := new(mocks.MockViewRepo)

	viewID := "view-1"
	teamID := "team-1"
	groupBys := []string{"status", "priority"}
	newName := "New View Name"

	editView := models.EditView{
		ID:       viewID,
		Name:     newName,
		TeamID:   teamID,
		GroupBys: groupBys,
	}

	data, _ := json.Marshal(editView)

	mockViewRepo.On("UpdateViewName", mock.Anything, viewID, newName).Return(nil)
	mockViewRepo.On("RemoveGroupByFromView", mock.Anything, viewID).Return(nil)
	mockViewRepo.On("RemoveIssueFromView", mock.Anything, viewID).Return(nil)

	mockViewRepo.On("AddGroupByToView", mock.Anything, mock.MatchedBy(func(p db.AddGroupByToViewParams) bool {
		return p.ViewID == viewID
	})).Return(nil)

	mockViewRepo.On("GetGroupedIssues", mock.Anything, teamID, groupBys).Return([]map[string]interface{}{
		{
			"status":   sql.NullString{Valid: true, String: "Open"},
			"priority": sql.NullString{Valid: true, String: "High"},
		},
	}, nil)

	mockViewRepo.On("ListIssuesByGroupFilters", mock.Anything, teamID, mock.MatchedBy(func(filter map[string]string) bool {
		return filter["status"] == "Open" && filter["priority"] == "High"
	})).Return([]db.Issue{
		{ID: "issue-1"},
		{ID: "issue-2"},
	}, nil)

	mockViewRepo.On("AddIssueToView", mock.Anything, mock.Anything).Return(nil)

	var broadcastCalled bool
	handler := &ws.WSHandler{
		ViewRepo: mockViewRepo,
		BroadcastFunc: func(data interface{}) {
			broadcastCalled = true
		},
	}

	mockConn := new(MockConn)
	mockConn.On("Locals", "userID").Return("user-1")

	handler.UpdateViewHandler(mockConn, data)

	mockViewRepo.AssertExpectations(t)
	if !broadcastCalled {
		t.Errorf("Expected broadcast to be called")
	}
}

