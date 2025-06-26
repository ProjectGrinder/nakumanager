package ws

import (
	"context"
	"encoding/json"
	"log"

	models "github.com/nack098/nakumanager/internal/models"
)

func (h *WSHandler) UpdateWorkspaceHandler(c ConnWithLocals, data json.RawMessage) {
	var payload models.EditWorkspace

	if err := json.Unmarshal(data, &payload); err != nil {
		log.Println("Invalid update_workspace payload:", err)
		return
	}

	ctx := context.Background()

	var event string
	var broadcastData map[string]interface{}

	if payload.Name != "" {
		err := h.WorkspaceRepo.RenameWorkspace(ctx, payload.WorkspaceID, payload.Name)
		if err != nil {
			log.Println("Failed to rename workspace:", err)
			return
		}

		event = "workspace_renamed"
		broadcastData = map[string]interface{}{
			"workspace_id": payload.WorkspaceID,
			"new_name":     payload.Name,
		}
	} else if payload.AddMember != "" {
		err := h.WorkspaceRepo.AddMemberToWorkspace(ctx, payload.WorkspaceID, payload.AddMember)
		if err != nil {
			log.Println("Failed to add member to workspace:", err)
			return
		}

		event = "workspace_member_added"
		broadcastData = map[string]interface{}{
			"workspace_id": payload.WorkspaceID,
			"user_id":      payload.AddMember,
		}
	} else if payload.RemoveMember != "" {
		err := h.WorkspaceRepo.RemoveMemberFromWorkspace(ctx, payload.WorkspaceID, payload.RemoveMember)
		if err != nil {
			log.Println("Failed to remove member from workspace:", err)
			return
		}

		event = "workspace_member_removed"
		broadcastData = map[string]interface{}{
			"workspace_id": payload.WorkspaceID,
			"user_id":      payload.RemoveMember,
		}
	} else {
		log.Println("No valid update action found in payload")
		return
	}

	log.Printf("Workspace %s updated with event %s", payload.WorkspaceID, event)

	h.Broadcast(map[string]interface{}{
		"event": event,
		"data":  broadcastData,
	})
}
