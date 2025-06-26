package ws_test

import (
	"encoding/json"
	"testing"

	"github.com/nack098/nakumanager/internal/db"
	models "github.com/nack098/nakumanager/internal/models"
	mocks "github.com/nack098/nakumanager/internal/routes/mock_repo"
	"github.com/nack098/nakumanager/internal/ws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockConn struct {
	mock.Mock
}

func (m *MockConn) Locals(key string, defaultValue ...interface{}) interface{} {
	args := m.Called(key)
	return args.Get(0)
}

func (m *MockConn) WriteMessage(messageType int, data []byte) error {
	args := m.Called(messageType, data)
	return args.Error(0)
}

func (m *MockConn) ReadMessage() (int, []byte, error) {
	args := m.Called()
	return args.Int(0), args.Get(1).([]byte), args.Error(2)
}

func (m *MockConn) Close() error {
	args := m.Called()
	return args.Error(0)
}

func strPtr(s string) *string {
	return &s
}

func TestUpdateProjectHandler_Success(t *testing.T) {
	mockProjectRepo := new(mocks.MockProjectRepo)
	handler := ws.NewWSHandler(nil, nil, mockProjectRepo, nil, nil, nil)

	mockConn := new(MockConn)
	mockConn.On("Locals", "userID").Return("user-123")

	projectID := "project-1"

	mockProjectRepo.On("IsProjectExists", mock.Anything, projectID).Return(true, nil)
	mockProjectRepo.On("GetProjectByID", mock.Anything, projectID).Return(db.Project{
		ID:        projectID,
		CreatedBy: "user-123",
	}, nil)
	mockProjectRepo.On("UpdateProjectName", mock.Anything, projectID, "New Project Name").Return(nil)
	mockProjectRepo.On("UpdateProjectLeader", mock.Anything, projectID, "leader-456").Return(nil)
	mockProjectRepo.On("UpdateProjectWorkspace", mock.Anything, projectID, "workspace-789").Return(nil)
	mockProjectRepo.On("AddMemberToProject", mock.Anything, db.AddMemberToProjectParams{
		ProjectID: projectID,
		UserID:    "member-001",
	}).Return(nil)
	mockProjectRepo.On("RemoveMemberFromProject", mock.Anything, db.RemoveMemberFromProjectParams{
		ProjectID: projectID,
		UserID:    "member-002",
	}).Return(nil)
	mockProjectRepo.On("UpdateProject", mock.Anything, mock.AnythingOfType("db.UpdateProjectParams")).Return(nil)

	var broadcastData map[string]interface{}
	handler.BroadcastFunc = func(data interface{}) {
		broadcastData = data.(map[string]interface{})
	}

	payload := models.EditProject{
		ID:           projectID,
		Name:         "New Project Name",
		LeaderID:     "leader-456",
		WorkspaceID:  "workspace-789",
		AddMember:    "member-001",
		RemoveMember: "member-002",
		Status:       strPtr("in_progress"),
		Priority:     strPtr("high"),
		Label:        strPtr("urgent"),
	}

	data, _ := json.Marshal(payload)
	raw := json.RawMessage(data)

	handler.UpdateProjectHandler(mockConn, raw)

	dataMap := broadcastData["data"].(map[string]interface{})

	assert.NotNil(t, broadcastData)
	assert.Equal(t, "project_updated", broadcastData["event"])
	assert.Equal(t, projectID, dataMap["id"])
	assert.Equal(t, "New Project Name", dataMap["name"])
	assert.Equal(t, "leader-456", dataMap["leader_id"])
	assert.Equal(t, "in_progress", *(dataMap["status"].(*string)))
	assert.Equal(t, "high", *(dataMap["priority"].(*string)))
	assert.Equal(t, "urgent", *(dataMap["label"].(*string)))

	mockConn.AssertExpectations(t)
	mockProjectRepo.AssertExpectations(t)
}
