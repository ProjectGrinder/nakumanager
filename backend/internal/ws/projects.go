package ws

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/nack098/nakumanager/internal/db"
	models "github.com/nack098/nakumanager/internal/models"
)

func ToNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func ToNullTime(t *time.Time) sql.NullTime {
	if t != nil {
		return sql.NullTime{Valid: true}
	}
	return sql.NullTime{Valid: false}
}

type ConnLocals interface {
	Locals(key string) interface{}
}

func (h *WSHandler) UpdateProjectHandler(c *websocket.Conn, data json.RawMessage) {
	log.Println("Received update_project event")

	var project models.EditProject
	if err := json.Unmarshal(data, &project); err != nil {
		log.Println("JSON unmarshal error:", err)
		return
	}

	log.Printf("Payload received: %+v\n", project)

	ctx := context.Background()

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		log.Println("unauthorized")
		return
	}
	log.Printf("User ID: %s\n", userID)

	_, err := h.ProjectRepo.IsProjectExists(ctx, project.ID)
	if err != nil {
		log.Println("Failed to check if project exists:", err)
		return
	}
	log.Println("Project exists, proceeding to update")

	aProject, err := h.ProjectRepo.GetProjectByID(ctx, project.ID)
	if err != nil {
		log.Println("Failed to check if user is owner of project:", err)
		return
	}
	if aProject.CreatedBy != userID {
		log.Println("User is not owner of project")
		return
	}
	log.Println("User is owner of project, proceeding to update")

	if project.Name != "" {
		err := h.ProjectRepo.UpdateProjectName(ctx, project.ID, project.Name)
		if err != nil {
			log.Println("Failed to rename project:", err)
			return
		}
		log.Printf("Project %s renamed to %s\n", project.ID, project.Name)
	}

	if project.LeaderID != "" {
		err := h.ProjectRepo.UpdateProjectLeader(ctx, project.ID, project.LeaderID)
		if err != nil {
			log.Println("Failed to update project leader:", err)
			return
		}
		log.Printf("Project %s leader updated to %s\n", project.ID, project.LeaderID)
	}

	if project.WorkspaceID != "" {
		err := h.ProjectRepo.UpdateProjectWorkspace(ctx, project.ID, project.WorkspaceID)
		if err != nil {
			log.Println("Failed to update project workspace:", err)
			return
		}
		log.Printf("Project %s workspace updated to %s\n", project.ID, project.WorkspaceID)
	}

	if project.AddMember != "" {
		err := h.ProjectRepo.AddMemberToProject(ctx, db.AddMemberToProjectParams{
			ProjectID: project.ID,
			UserID:    project.AddMember,
		})
		if err != nil {
			log.Println("Failed to add member to project:", err)
			return
		}
		log.Printf("Member %s added to project %s\n", project.AddMember, project.ID)
	}

	if project.RemoveMember != "" {
		err := h.ProjectRepo.RemoveMemberFromProject(ctx, db.RemoveMemberFromProjectParams{
			ProjectID: project.ID,
			UserID:    project.RemoveMember,
		})
		if err != nil {
			log.Println("Failed to remove member from project:", err)
			return
		}
		log.Printf("Member %s removed from project %s\n", project.RemoveMember, project.ID)
	}

	err = h.ProjectRepo.UpdateProject(ctx, db.UpdateProjectParams{
		ID:        project.ID,
		Status:    ToNullString(project.Status),
		Priority:  ToNullString(project.Priority),
		Label:     ToNullString(project.Label),
		StartDate: ToNullTime(project.StartDate),
		EndDate:   ToNullTime(project.EndDate),
	})
	if err != nil {
		log.Println("Failed to update project:", err)
		return
	}

	log.Printf("Project %s updated successfully\n", project.ID)

	h.Broadcast(map[string]interface{}{
		"event": "project_updated",
		"data": map[string]interface{}{
			"id":         project.ID,
			"name":       project.Name,
			"leader_id":  project.LeaderID,
			"status":     project.Status,
			"priority":   project.Priority,
			"label":      project.Label,
			"start_date": project.StartDate,
			"end_date":   project.EndDate,
		},
	})
}
