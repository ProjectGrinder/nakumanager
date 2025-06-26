package ws

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/nack098/nakumanager/internal/db"
	models "github.com/nack098/nakumanager/internal/models"
)

func (h *WSHandler) UpdateTeamHandler(c *websocket.Conn, data json.RawMessage) {
	var payload models.EditTeam

	if err := json.Unmarshal(data, &payload); err != nil {
		log.Println("Invalid update_team payload:", err)
		return
	}

	ctx := context.Background()

	team, err := h.TeamRepo.GetTeamByID(ctx, payload.TeamID)
	if err != nil {
		log.Println("Team not found")
		return
	}

	var event string
	var broadcastData map[string]interface{}

	if payload.Name != "" {
		err := h.TeamRepo.RenameTeam(ctx, payload.TeamID, payload.Name)
		if err != nil {
			log.Println("Failed to rename team:", err)
			return
		}

		event = "team_renamed"
		broadcastData = map[string]interface{}{
			"team_id":  payload.TeamID,
			"new_name": payload.Name,
		}
	}
	if payload.AddMember != "" {
		err := h.TeamRepo.AddMemberToTeam(ctx, db.AddMemberToTeamParams{
			UserID: payload.AddMember,
			TeamID: payload.TeamID,
		})
		if err != nil {
			log.Println("Failed to add member to workspace:", err)
			return
		}

		event = "team_member_added"
		broadcastData = map[string]interface{}{
			"workspace_id": payload.TeamID,
			"user_id":      payload.AddMember,
		}
	}
	if payload.RemoveMember != "" {
		err := h.TeamRepo.RemoveMemberFromTeam(ctx, db.RemoveMemberFromTeamParams{
			TeamID: payload.TeamID,
			UserID: payload.RemoveMember,
		})
		if err != nil {
			log.Println("Failed to remove member from team:", err)
			return
		}
		if team.LeaderID.Valid && payload.RemoveMember == team.LeaderID.String {
			err := h.TeamRepo.SetLeaderToTeam(ctx, payload.TeamID, "")
			if err != nil {
				log.Println("Failed to set leader to team:", err)
				return
			}
		}

		event = "team_member_removed"
		broadcastData = map[string]interface{}{
			"workspace_id": payload.TeamID,
			"user_id":      payload.RemoveMember,
		}
	}
	if payload.Leader != "" {
		log.Println("Leader: ", payload.Leader)
		err := h.TeamRepo.SetLeaderToTeam(ctx, payload.TeamID, payload.Leader)
		if err != nil {
			log.Println("Failed to set leader to team:", err)
			return
		}

		event = "team_leader_set"
		broadcastData = map[string]interface{}{
			"workspace_id": payload.TeamID,
			"user_id":      payload.Leader,
		}

	}

	if event == "" {
		log.Println("No valid update action found in payload")
		return
	}

	log.Printf("Team %s updated with event %s", payload.TeamID, event)

	h.Broadcast(map[string]interface{}{
		"event": event,
		"data":  broadcastData,
	})

}
