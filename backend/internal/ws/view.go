package ws

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"

	"github.com/gofiber/contrib/websocket"
	"github.com/nack098/nakumanager/internal/db"
	models "github.com/nack098/nakumanager/internal/models"
)

func (h *WSHandler) UpdateViewHandler(c *websocket.Conn, data json.RawMessage) {
	var view models.EditView

	if err := json.Unmarshal(data, &view); err != nil {
		log.Println("json error:", err)
		return
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		log.Println("Unauthorized: missing userID in connection locals")
		return
	}

	ctx := context.Background()
	log.Printf("UpdateViewHandler called by user: %s for view ID: %s", userID, view.ID)

	if view.Name != "" {
		log.Printf("Updating view name to: %s", view.Name)
		err := h.ViewRepo.UpdateViewName(ctx, view.ID, view.Name)
		if err != nil {
			log.Println("Failed to rename view:", err)
			return
		}
	}

	if len(view.GroupBys) < 1 || len(view.GroupBys) > 2 {
		log.Printf("Invalid group_bys length: %d, must be 1 or 2", len(view.GroupBys))
		return
	} else {
		log.Printf("Removing old group_bys for view ID: %s", view.ID)
		err := h.ViewRepo.RemoveGroupByFromView(ctx, view.ID)
		if err != nil {
			log.Printf("Failed to delete group_by: %v", err)
			return
		}

		log.Printf("Removing old issues from view ID: %s", view.ID)
		err = h.ViewRepo.RemoveIssueFromView(ctx, view.ID)
		if err != nil {
			log.Printf("Failed to delete issues: %v", err)
			return
		}

		for _, g := range view.GroupBys {
			log.Printf("Adding group_by '%s' to view ID: %s", g, view.ID)
			err := h.ViewRepo.AddGroupByToView(ctx, db.AddGroupByToViewParams{
				ViewID:  view.ID,
				GroupBy: g,
			})
			if err != nil {
				log.Printf("Failed to add group_by: %s -> %v", g, err)
				return
			}
		}

		log.Printf("Getting grouped issues by %v for team ID: %s", view.GroupBys, view.TeamID)
		groupedData, err := h.ViewRepo.GetGroupedIssues(ctx, view.TeamID, view.GroupBys)
		if err != nil {
			log.Printf("Failed to group issues by %+v: %v", view.GroupBys, err)
			return
		}

		issueSet := make(map[string]struct{})

		for _, row := range groupedData {
			filter := make(map[string]string)
			skip := false

			for _, col := range view.GroupBys {
				if val, ok := row[col]; ok {
					if v, ok := val.(sql.NullString); ok {
						if v.Valid {
							filter[col] = v.String
						} else {
							skip = true
							break
						}
					}
				}
			}

			if skip {
				continue
			}

			log.Printf("Listing issues with filter: %+v", filter)
			matchingIssues, err := h.ViewRepo.ListIssuesByGroupFilters(ctx, view.TeamID, filter)
			if err != nil {
				log.Printf("Failed to list issues by group %+v: %v", filter, err)
				return
			}

			for _, issue := range matchingIssues {
				issueSet[issue.ID] = struct{}{}
			}
		}

		for issueID := range issueSet {
			log.Printf("Adding issue %s to view %s", issueID, view.ID)
			err := h.ViewRepo.AddIssueToView(ctx, db.AddIssueToViewParams{
				ViewID:  view.ID,
				IssueID: issueID,
			})
			if err != nil {
				log.Printf("Failed to add issue %s to view: %v", issueID, err)
				return
			}
		}
	}

	log.Printf("Broadcasting view update for view ID: %s", view.ID)
	h.Broadcast(map[string]interface{}{
		"view_id":   view.ID,
		"name":      view.Name,
		"group_bys": view.GroupBys,
	})
}
