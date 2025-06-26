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
	"github.com/stretchr/testify/require"
)

func TestUpdateTeamHandler_Success(t *testing.T) {
	mockTeamRepo := new(mocks.MockTeamRepo)

	mockTeamRepo.On("GetTeamByID", mock.Anything, "team-123").Return(db.Team{
		ID: "team-123",
	}, nil)

	mockTeamRepo.On("RenameTeam", mock.Anything, "team-123", "New Team Name").Return(nil)

	mockTeamRepo.On("AddMemberToTeam", mock.Anything, db.AddMemberToTeamParams{
		UserID: "user-456",
		TeamID: "team-123",
	}).Return(nil)

	mockTeamRepo.On("RemoveMemberFromTeam", mock.Anything, db.RemoveMemberFromTeamParams{
		TeamID: "team-123",
		UserID: "user-456",
	}).Return(nil)

	mockTeamRepo.On("SetLeaderToTeam", mock.Anything, "team-123", "user-789").Return(nil)

	handler := ws.NewWSHandler(nil, mockTeamRepo, nil, nil, nil, nil)

	var broadcasts []map[string]interface{}
	handler.BroadcastFunc = func(msg interface{}) {
		broadcasts = append(broadcasts, msg.(map[string]interface{}))
	}

	payloadRename := models.EditTeam{
		TeamID: "team-123",
		Name:   "New Team Name",
	}
	dataRename, _ := json.Marshal(payloadRename)
	handler.UpdateTeamHandler(nil, dataRename)

	payloadAddMember := models.EditTeam{
		TeamID:    "team-123",
		AddMember: "user-456",
	}
	dataAddMember, _ := json.Marshal(payloadAddMember)
	handler.UpdateTeamHandler(nil, dataAddMember)

	payloadRemoveMember := models.EditTeam{
		TeamID:       "team-123",
		RemoveMember: "user-456",
	}
	dataRemoveMember, _ := json.Marshal(payloadRemoveMember)
	handler.UpdateTeamHandler(nil, dataRemoveMember)

	payloadSetLeader := models.EditTeam{
		TeamID: "team-123",
		Leader: "user-789",
	}
	dataSetLeader, _ := json.Marshal(payloadSetLeader)
	handler.UpdateTeamHandler(nil, dataSetLeader)


	require.Len(t, broadcasts, 4)

	assert.Equal(t, "team_renamed", broadcasts[0]["event"])
	assert.Equal(t, "team-123", broadcasts[0]["data"].(map[string]interface{})["team_id"])
	assert.Equal(t, "New Team Name", broadcasts[0]["data"].(map[string]interface{})["new_name"])

	assert.Equal(t, "team_member_added", broadcasts[1]["event"])
	assert.Equal(t, "team-123", broadcasts[1]["data"].(map[string]interface{})["workspace_id"])
	assert.Equal(t, "user-456", broadcasts[1]["data"].(map[string]interface{})["user_id"])

	assert.Equal(t, "team_member_removed", broadcasts[2]["event"])
	assert.Equal(t, "team-123", broadcasts[2]["data"].(map[string]interface{})["workspace_id"])
	assert.Equal(t, "user-456", broadcasts[2]["data"].(map[string]interface{})["user_id"])

	assert.Equal(t, "team_leader_set", broadcasts[3]["event"])
	assert.Equal(t, "team-123", broadcasts[3]["data"].(map[string]interface{})["workspace_id"])
	assert.Equal(t, "user-789", broadcasts[3]["data"].(map[string]interface{})["user_id"])

	mockTeamRepo.AssertExpectations(t)
}
