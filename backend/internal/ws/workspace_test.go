package ws_test

import (
	"encoding/json"
	"testing"

	models "github.com/nack098/nakumanager/internal/models"
	mocks "github.com/nack098/nakumanager/internal/routes/mock_repo"
	"github.com/nack098/nakumanager/internal/ws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)


func TestUpdateWorkspaceHandler_RenameWorkspace(t *testing.T) {
	mockRepo := new(mocks.MockWorkspaceRepo)
	mockRepo.On("RenameWorkspace", mock.Anything, "ws-123", "New Name").Return(nil)

	handler := ws.NewWSHandler(mockRepo, nil, nil, nil, nil, nil)

	var broadcastCalled bool
	var broadcastData map[string]interface{}
	handler.BroadcastFunc = func(msg interface{}) {
		broadcastCalled = true
		broadcastData = msg.(map[string]interface{})
	}

	payload := models.EditWorkspace{
		WorkspaceID: "ws-123",
		Name:        "New Name",
	}
	data, _ := json.Marshal(payload)

	handler.UpdateWorkspaceHandler(nil, data)

	assert.True(t, broadcastCalled, "Broadcast should be called")
	assert.Equal(t, "workspace_renamed", broadcastData["event"])
	dataMap := broadcastData["data"].(map[string]interface{})
	assert.Equal(t, "ws-123", dataMap["workspace_id"])
	assert.Equal(t, "New Name", dataMap["new_name"])

	mockRepo.AssertExpectations(t)
}

func TestUpdateWorkspaceHandler_AddMemberToWorkspace(t *testing.T) {
	mockRepo := new(mocks.MockWorkspaceRepo)
	mockRepo.On("AddMemberToWorkspace", mock.Anything, "ws-123", "user-456").Return(nil)

	handler := ws.NewWSHandler(mockRepo, nil, nil, nil, nil, nil)

	var broadcastCalled bool
	var broadcastData map[string]interface{}
	handler.BroadcastFunc = func(msg interface{}) {
		broadcastCalled = true
		broadcastData = msg.(map[string]interface{})
	}

	payload := models.EditWorkspace{
		WorkspaceID: "ws-123",
		AddMember:   "user-456",
	}
	data, _ := json.Marshal(payload)

	handler.UpdateWorkspaceHandler(nil, data)

	assert.True(t, broadcastCalled, "Broadcast should be called")
	assert.Equal(t, "workspace_member_added", broadcastData["event"])
	dataMap := broadcastData["data"].(map[string]interface{})
	assert.Equal(t, "ws-123", dataMap["workspace_id"])
	assert.Equal(t, "user-456", dataMap["user_id"])

	mockRepo.AssertExpectations(t)
}

func TestUpdateWorkspaceHandler_RemoveMemberFromWorkspace(t *testing.T) {
	mockRepo := new(mocks.MockWorkspaceRepo)
	mockRepo.On("RemoveMemberFromWorkspace", mock.Anything, "ws-123", "user-456").Return(nil)

	handler := ws.NewWSHandler(mockRepo, nil, nil, nil, nil, nil)

	var broadcastCalled bool
	var broadcastData map[string]interface{}
	handler.BroadcastFunc = func(msg interface{}) {
		broadcastCalled = true
		broadcastData = msg.(map[string]interface{})
	}

	payload := models.EditWorkspace{
		WorkspaceID:  "ws-123",
		RemoveMember: "user-456",
	}
	data, _ := json.Marshal(payload)

	handler.UpdateWorkspaceHandler(nil, data)

	assert.True(t, broadcastCalled, "Broadcast should be called")
	assert.Equal(t, "workspace_member_removed", broadcastData["event"])
	dataMap := broadcastData["data"].(map[string]interface{})
	assert.Equal(t, "ws-123", dataMap["workspace_id"])
	assert.Equal(t, "user-456", dataMap["user_id"])

	mockRepo.AssertExpectations(t)
}
