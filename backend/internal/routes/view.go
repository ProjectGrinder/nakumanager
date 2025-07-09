package routes

import (
	"database/sql"
	"fmt"
	"log"

	models "github.com/nack098/nakumanager/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/nack098/nakumanager/internal/db"
	"github.com/nack098/nakumanager/internal/repositories"
)

type ViewHandler struct {
	DB   *sql.DB
	Repo repositories.ViewRepository
}

func NewViewHandler(db *sql.DB, repo repositories.ViewRepository) *ViewHandler {
	return &ViewHandler{
		DB:   db,
		Repo: repo}
}

func safeString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return "null"
}

func isSafeColumn(col string) bool {
	allowed := map[string]bool{
		"status":     true,
		"priority":   true,
		"assignee":   true,
		"project_id": true,
		"label":      true,
		"team_id":    true,
		"end_date":   true,
	}
	return allowed[col]
}

func (h *ViewHandler) CreateView(c *fiber.Ctx) error {
	var req models.CreateView
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if req.Name == "" || req.TeamID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "name and team_id are required"})
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	req.ID = uuid.New().String()
	ctx := c.Context()

	if err := h.Repo.CreateView(ctx, db.CreateViewParams{
		ID:        req.ID,
		Name:      req.Name,
		CreatedBy: userID,
		TeamID:    req.TeamID,
	}); err != nil {
		log.Printf("Failed to create view: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create view"})
	}

	validGroupBys := map[string]bool{
		"status": true, "priority": true, "project_id": true,
		"label": true, "assignee": true, "team_id": true, "end_date": true,
	}
	for _, g := range req.GroupBys {
		if !validGroupBys[g] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid group_by value: %s", g),
			})
		}
		if err := h.Repo.AddGroupByToView(ctx, db.AddGroupByToViewParams{
			ViewID:  req.ID,
			GroupBy: g,
		}); err != nil {
			log.Printf("Failed to add group_by %s: %v", g, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to add group_by: %s", g),
			})
		}
	}

	query, err := buildGroupByQuery("issues", req.GroupBys)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	rows, err := h.DB.QueryContext(ctx, query, req.TeamID)
	if err != nil {
		log.Printf("Failed to query grouped issues: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to group issues"})
	}
	defer rows.Close()

	issueSet := make(map[string]struct{})
	for rows.Next() {
		var issueID string
		if err := rows.Scan(&issueID); err == nil {
			issueSet[issueID] = struct{}{}
		}
	}
	if err := rows.Err(); err != nil {
		log.Printf("Row error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read grouped issues"})
	}

	for issueID := range issueSet {
		if err := h.Repo.AddIssueToView(ctx, db.AddIssueToViewParams{
			ViewID:  req.ID,
			IssueID: issueID,
		}); err != nil {
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
	var req models.ViewGroupBy
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if len(req.GroupBys) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "group_by is required"})
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	views, err := h.Repo.GetViewsByGroupBys(c.Context(), req.GroupBys)
	if err != nil {
		log.Printf("Failed to get views by group_by: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get views by group_by"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Success", "views": views})
}

func (h *ViewHandler) GetViewByTeamID(c *fiber.Ctx) error {
	log.Println("GetViewByTeamID")
	teamID := c.Params("id")
	if teamID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body or missing team_id",
		})
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	views, err := h.Repo.ListViewByTeamID(c.Context(), teamID)
	if err != nil {
		log.Printf("Failed to get views by user ID: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get views",
		})
	}

	return c.Status(fiber.StatusOK).JSON(views)
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

func (h *ViewHandler) UpdateView(c *fiber.Ctx) error {
	viewID := c.Params("id")
	if viewID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body or missing view_id",
		})
	}

	var req models.UpdateViewRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	teamID, err := h.Repo.GetTeamIDByViewID(c.Context(), viewID)
	if err != nil {
		log.Printf("Failed to get team_id: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get team_id"})
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	log.Printf("UserID: %s is updating ViewID: %s", userID, viewID)

	view, err := h.Repo.GetViewByID(c.Context(), viewID)
	if err != nil {
		log.Printf("Failed to get view by ID: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get view"})
	}
	if len(view) == 0 {
		log.Printf("View not found: %s", viewID)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "View not found"})
	}

	ctx := c.Context()
	tx, err := h.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to begin transaction"})
	}

	defer func() {
		if err != nil {
			log.Println("Rolling back transaction...")
			tx.Rollback()
		} else {
			log.Println("Committing transaction...")
			tx.Commit()
		}
	}()

	if req.Name != "" {
		log.Printf("Updating view name to: %s", req.Name)
		err = h.Repo.UpdateViewName(ctx, viewID, req.Name)
		if err != nil {
			log.Printf("Failed to update view name: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update view name"})
		}
	}

	if req.TeamID != "" {
		log.Printf("Updating view team_id to: %s", req.TeamID)
		err = h.Repo.UpdateViewTeamID(ctx, viewID, req.TeamID)
		if err != nil {
			log.Printf("Failed to update view team_id: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update view team_id"})
		}
	}

	if req.GroupBys != nil {
		log.Printf("Updating group_bys to: %v", req.GroupBys)

		if err = h.Repo.RemoveGroupByFromView(ctx, viewID); err != nil {
			log.Printf("Failed to remove old group_bys: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove old group_bys"})
		}
		if err = h.Repo.RemoveIssueFromView(ctx, viewID); err != nil {
			log.Printf("Failed to remove old issues: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove old issues"})
		}

		for _, gb := range req.GroupBys {
			log.Printf("Adding group_by: %s", gb)
			err = h.Repo.AddGroupByToView(ctx, db.AddGroupByToViewParams{
				ViewID:  viewID,
				GroupBy: gb,
			})
			if err != nil {
				log.Printf("Failed to add group_by %s: %v", gb, err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to add group_by"})
			}
		}

		query, err := buildGroupByQuery("issues", req.GroupBys)
		if err != nil {
			log.Printf("Failed to build query: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		log.Printf("Using team_id: %s", teamID)

		rows, err := tx.QueryContext(ctx, query, teamID)
		if err != nil {
			log.Printf("Failed to query grouped issues: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to group issues"})
		}
		defer rows.Close()

		issueSet := make(map[string]struct{})
		found := false
		for rows.Next() {
			var issueID string
			if err := rows.Scan(&issueID); err == nil {
				log.Printf("Found issue ID: %s", issueID)
				issueSet[issueID] = struct{}{}
				found = true
			} else {
				log.Printf("Row scan error: %v", err)
			}
		}
		if !found {
			log.Println("No issues found by group_by query")
		}
		if err := rows.Err(); err != nil {
			log.Printf("Row error: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read grouped issues"})
		}

		for issueID := range issueSet {
			log.Printf("Adding issue %s to view %s", issueID, viewID)
			if err := h.Repo.AddIssueToViewTx(ctx, tx, db.AddIssueToViewParams{
				ViewID:  viewID,
				IssueID: issueID,
			}); err != nil {
				log.Printf("Failed to add issue %s to view: %v", issueID, err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": fmt.Sprintf("Failed to add issue %s to view", issueID),
				})
			}
		}
	}

	log.Println("View updated successfully")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "View updated successfully"})
}
