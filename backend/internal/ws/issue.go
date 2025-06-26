package ws

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nack098/nakumanager/internal/db"
	models "github.com/nack098/nakumanager/internal/models"
)

func (h *WSHandler) UpdateIssueHandler(c ConnWithLocals, data json.RawMessage) {
	log.Println("Received update_issue event")

	var issue models.EditIssue
	if err := json.Unmarshal(data, &issue); err != nil {
		log.Println("Failed to unmarshal update issue payload:", err)
		return
	}
	log.Printf("Payload received: %+v\n", issue)

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		log.Println("Unauthorized: missing userID in context")
		return
	}

	ctx := context.Background()

	existing, err := h.IssueRepo.GetIssueByID(ctx, issue.ID)
	if err != nil {
		log.Printf("Issue not found: %v", err)
		return
	}

	if issue.Assignee == nil && existing.OwnerID != userID {
		log.Println("Unauthorized update_issue: not owner and no assignee")
		return
	}

	if issue.Title == nil || issue.Status == nil || issue.TeamID == nil || issue.OwnerID == nil {
		log.Println("Missing required fields (Title, Status, TeamID, or OwnerID cannot be nil)")
		return
	}

	isMember, err := h.TeamRepo.IsMemberInTeam(ctx, *issue.TeamID, userID)
	if err != nil {
		log.Println("Failed to check team membership:", err)
		return
	}
	if !isMember {
		log.Println("User is not a member of the team, update denied")
		return
	}

	if issue.Assignee != nil {
		log.Printf("Checking assignee existence: %s", *issue.Assignee)
		_, err := h.UserRepo.GetUserByID(ctx, *issue.Assignee)
		if err != nil {
			log.Println("Assignee does not exist:", err)
			return
		}

		log.Printf("Adding assignee to issue: %s", *issue.Assignee)
		if err := h.IssueRepo.AddAssigneeToIssue(ctx, db.AddAssigneeToIssueParams{
			IssueID: issue.ID,
			UserID:  *issue.Assignee,
		}); err != nil {
			log.Println("Failed to add assignee to issue:", err)
			return
		}
	}

	if issue.StartDate != nil && issue.EndDate != nil && issue.EndDate.Before(*issue.StartDate) {
		log.Println("Invalid dates: EndDate is before StartDate")
		return
	}

	log.Printf("Updating issue: ID=%s, Title=%s, Status=%s, TeamID=%s, OwnerID=%s",
		issue.ID, *issue.Title, *issue.Status, *issue.TeamID, *issue.OwnerID)

	log.Printf("Additional fields - Content: %v, Priority: %v, ProjectID: %v, StartDate: %v, EndDate: %v, Label: %v",
		ToNullString(issue.Content), ToNullString(issue.Priority), ToNullString(issue.ProjectID),
		ToNullTime(issue.StartDate), ToNullTime(issue.EndDate), ToNullString(issue.Label))

	err = h.IssueRepo.UpdateIssue(ctx, db.UpdateIssueParams{
		ID:        issue.ID,
		Title:     *issue.Title,
		Content:   ToNullString(issue.Content),
		Priority:  ToNullString(issue.Priority),
		Status:    *issue.Status,
		ProjectID: ToNullString(issue.ProjectID),
		TeamID:    *issue.TeamID,
		StartDate: ToNullTime(issue.StartDate),
		EndDate:   ToNullTime(issue.EndDate),
		Label:     ToNullString(issue.Label),
		OwnerID:   *issue.OwnerID,
	})
	if err != nil {
		log.Println("Failed to update issue:", err)
		return
	}

	log.Println("Issue updated successfully:", issue.ID)

	broadcastData := map[string]interface{}{
		"event": "issue_updated",
		"data":  issue,
	}

	log.Printf("Broadcasting update for issue %s\n", issue.ID)
	h.Broadcast(broadcastData)
}
