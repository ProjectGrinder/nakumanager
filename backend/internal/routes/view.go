package routes

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	models "github.com/nack098/nakumanager/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nack098/nakumanager/internal/db"
	"github.com/nack098/nakumanager/internal/repositories"
)

type ViewHandler struct {
	Repo repositories.ViewRepository
}

func NewViewHandler(repo repositories.ViewRepository) *ViewHandler {
	return &ViewHandler{Repo: repo}
}

func safeString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return "null"
}


func (h *ViewHandler) CreateView(c *fiber.Ctx) error {
	var req models.ViewCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Name == "" || req.TeamID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "name and team_id are required",
		})
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	req.ID = uuid.New().String()


	err := h.Repo.CreateView(c.Context(), db.CreateViewParams{
		ID:        req.ID,
		Name:      req.Name,
		CreatedBy: userID,
		TeamID:    req.TeamID,
	})
	if err != nil {
		log.Printf("Failed to create view: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create view",
		})
	}

	validGroupBys := map[string]bool{
		"status": true, "assignee": true, "priority": true,
		"project_id": true, "label": true, "team_id": true, "end_date": true,
	}

	for _, g := range req.GroupBys {
		if !validGroupBys[g] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid group_by value: %s", g),
			})
		}
	}

	for _, g := range req.GroupBys {
		err := h.Repo.AddGroupByToView(c.Context(), db.AddGroupByToViewParams{
			ViewID:  req.ID,
			GroupBy: g,
		})
		if err != nil {
			log.Printf("Failed to add group_by: %s -> %v", g, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to add group_by: %s", g),
			})
		}
	}

	groupedData, err := h.Repo.GetGroupedIssues(c.Context(), req.TeamID, req.GroupBys)
	if err != nil {
		log.Printf("Failed to group issues by %+v: %v", req.GroupBys, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to group issues",
		})
	}

	issueSet := make(map[string]struct{})

	for _, row := range groupedData {
		filter := make(map[string]string)
		for _, col := range req.GroupBys {
			if val, ok := row[col]; ok {
				if v, ok := val.(sql.NullString); ok && v.Valid {
					filter[col] = v.String
				}
			}
		}

		matchingIssues, err := h.Repo.ListIssuesByGroupFilters(c.Context(), req.TeamID, filter)
		if err != nil {
			log.Printf("Failed to list issues by group %+v: %v", filter, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to find issues by group %v", filter),
			})
		}

		for _, issue := range matchingIssues {
			issueSet[issue.ID] = struct{}{}
		}
	}

	for issueID := range issueSet {
		err := h.Repo.AddIssueToView(c.Context(), db.AddIssueToViewParams{
			ViewID:  req.ID,
			IssueID: issueID,
		})
		if err != nil {
			log.Printf("Failed to add issue %s to view: %v", issueID, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to add issue %s to view", issueID),
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "View created successfully",
		"id":      req.ID,
	})
}

func (h *ViewHandler) GetViewsByGroupBy(c *fiber.Ctx) error {
	var viewGroupBy models.ViewGroupBy
	if err := c.BodyParser(&viewGroupBy); err != nil || viewGroupBy.ViewID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body or missing view_id",
		})
	}

	groupByCols, err := h.Repo.ListGroupByViewID(c.UserContext(), viewGroupBy.ViewID)
	if err != nil {
		log.Printf("Failed to get group_by: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get group_by settings",
		})
	}

	issues, err := h.Repo.ListIssuesByViewID(c.UserContext(), viewGroupBy.ViewID)
	if err != nil {
		log.Printf("Failed to get issues: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get issues",
		})
	}

	for _, issue := range issues {
		log.Printf("[DEBUG] Issue ID = %s, Status = %s, Priority = %+v, ProjectID = %+v", issue.ID, issue.Status, issue.Priority, issue.ProjectID)
	}

	groupMap := make(map[string][]db.Issue)
	for _, issue := range issues {
		var keyParts []string
		for _, col := range groupByCols {
			switch col {
			case "status":
				keyParts = append(keyParts, issue.Status)
			case "priority":
				keyParts = append(keyParts, safeString(issue.Priority))
			case "project_id":
				keyParts = append(keyParts, safeString(issue.ProjectID))
			case "assignee":
				keyParts = append(keyParts, safeString(issue.Assignee))
			case "label":
				keyParts = append(keyParts, safeString(issue.Label))
			case "end_date":
				if !issue.EndDate.Time.IsZero() {
					keyParts = append(keyParts, issue.EndDate.Time.Format("2006-01-02"))
				} else {
					keyParts = append(keyParts, "null")
				}
			default:
				keyParts = append(keyParts, "unknown")
			}
		}

		groupKey := strings.Join(keyParts, "|")

		groupMap[groupKey] = append(groupMap[groupKey], issue)
	}

	var results []map[string]interface{}
	for key, issues := range groupMap {
		group := map[string]interface{}{}
		keyParts := strings.Split(key, "|")
		for i, col := range groupByCols {
			group[col] = keyParts[i]
		}
		group["issues"] = issues
		results = append(results, group)
	}

	log.Printf("[DEBUG] Total groups: %d", len(results))

	return c.Status(fiber.StatusOK).JSON(results)
}

func (h *ViewHandler) DeleteView(c *fiber.Ctx) error {
	viewID := c.Params("id")
	if viewID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body or missing view_id",
		})
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	view, err := h.Repo.GetViewByID(c.Context(), viewID)
	if err != nil {
		log.Printf("Failed to get view by ID: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get view",
		})
	}

	log.Printf("[DEBUG] View details: %+v", view)

	if len(view) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "View not found",
		})
	}

	if view[0].CreatedBy != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "You don't have permission to delete this view",
		})
	}

	err = h.Repo.DeleteView(c.Context(), viewID)
	if err != nil {
		log.Printf("Failed to delete view: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete view",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "View deleted successfully",
	})
}
